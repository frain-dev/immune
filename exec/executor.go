package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/database"
	"github.com/frain-dev/immune/url"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Executor is used to execute tests
type Executor struct {
	callbackIDLocation     string
	baseURL                string
	maxCallbackWaitSeconds uint
	idFn                   func() string
	client                 *http.Client
	dbTruncator            database.Truncator
	vm                     *immune.VariableMap
	s                      immune.CallbackServer
}

func NewExecutor(
	s immune.CallbackServer,
	client *http.Client,
	vm *immune.VariableMap,
	maxCallbackWaitSeconds uint,
	baseURL string,
	callbackIDLocation string,
	dbTruncator database.Truncator, idFn func() string) *Executor {
	return &Executor{
		s:                      s,
		vm:                     vm,
		idFn:                   idFn,
		client:                 client,
		baseURL:                baseURL,
		dbTruncator:            dbTruncator,
		callbackIDLocation:     callbackIDLocation,
		maxCallbackWaitSeconds: maxCallbackWaitSeconds,
	}
}

// ExecuteSetupTestCase executes setup test cases
func (ex *Executor) ExecuteSetupTestCase(ctx context.Context, setupTC *immune.SetupTestCase) error {
	u, err := url.Parse(fmt.Sprintf("%s%s", ex.baseURL, setupTC.Endpoint))
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %s: failed to parse url", setupTC.Name)
	}

	result, err := u.ProcessWithVariableMap(ex.vm)
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %s: failed to process parsed url with variable map", setupTC.Name)
	}

	r := &request{
		contentType: "application/json",
		url:         result,
		body:        setupTC.RequestBody,
		method:      setupTC.HTTPMethod,
	}

	err = r.processWithVariableMap(ex.vm)
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %s: failed to process request body with variable map", setupTC.Name)
	}
	//log.Infof("setup_test_case %s: request body: %s", setupTC.Name, pretty.Sprint(r.body))

	resp, err := ex.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	if setupTC.StatusCode != resp.statusCode {
		return errors.Errorf("setup_test_case %s: wants status code %d but got status code %d, response body: %s", setupTC.Name, setupTC.StatusCode, resp.statusCode, resp.body.String())
	}

	if setupTC.ResponseBody {
		if resp.body.Len() == 0 {
			return errors.Errorf("setup_test_case %s: wants response body but got no response body", setupTC.Name)
		}

		m := immune.M{}
		err = resp.Decode(&m)
		if err != nil {
			return errors.Wrapf(err, "setup_test_case %s: failed to decode response body: response body: %s", setupTC.Name, resp.body.String())
		}

		if setupTC.StoreResponseVariables != nil {
			err = ex.vm.ProcessResponse(ctx, setupTC.StoreResponseVariables, m)
			if err != nil {
				return errors.Wrapf(err, "setup_test_case %s: failed to process response body: response body: %s", setupTC.Name, resp.body.String())
			}
		}
	} else {
		if resp.body.Len() > 0 {
			return errors.Errorf("setup_test_case %s: does not want a response body but got a response body: '%s'", setupTC.Name, resp.body.String())
		}
	}

	return nil
}

// ExecuteTestCase executes test cases, it waits for callback if necessary
func (ex *Executor) ExecuteTestCase(ctx context.Context, tc *immune.TestCase) error {
	u, err := url.Parse(fmt.Sprintf("%s%s", ex.baseURL, tc.Endpoint))
	if err != nil {
		return errors.Wrapf(err, "test_case %s: failed to parse url", tc.Name)
	}

	result, err := u.ProcessWithVariableMap(ex.vm)
	if err != nil {
		return errors.Wrapf(err, "test_case %s: failed to process parsed url with variable map", tc.Name)
	}

	var uid string
	if tc.Callback.Enabled {
		uid = ex.idFn()
		err = immune.InjectCallbackID(ex.callbackIDLocation, uid, tc.RequestBody)
		if err != nil {
			return errors.Wrapf(err, "test_case %s: failed to inject callback id into request body", tc.Name)
		}
	}

	r := &request{
		contentType: "application/json",
		body:        tc.RequestBody,
		url:         result,
		method:      tc.HTTPMethod,
	}

	err = r.processWithVariableMap(ex.vm)
	if err != nil {
		return errors.Wrapf(err, "test_case %s: failed to process request body with variable map", tc.Name)
	}

	resp, err := ex.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	if tc.StatusCode != resp.statusCode {
		return errors.Errorf("test_case %s: wants status code %d but got status code %d", tc.Name, tc.StatusCode, resp.statusCode)
	}

	if tc.ResponseBody {
		if resp.body.Len() == 0 {
			return errors.Errorf("test_case %s: wants response body but got no response body: status_code: %d", tc.Name, resp.statusCode)
		}

		m := immune.M{}
		err = resp.Decode(&m)
		if err != nil {
			return errors.Wrapf(err, "test_case %s: failed to decode response body: %s", tc.Name, string(resp.buf))
		}

	} else {
		if resp.body.Len() > 0 {
			return errors.Errorf("test_case %s: does not want a response body but got a response body: '%s'", tc.Name, resp.body.String())
		}
	}

	if tc.Callback.Enabled {
		cctx, cancel := context.WithTimeout(context.Background(), time.Duration(ex.maxCallbackWaitSeconds)*time.Second)
		defer cancel()

		signalChan := make(chan *immune.Signal, 1)

		for i := uint(1); i <= tc.Callback.Times; i++ {
			ex.s.ReceiveCallback(signalChan)

			select {
			case <-cctx.Done():
				log.Infof("succesfully received %d callbacks for test_case %s before max callback wait seconds elapsed", i-1, tc.Name)
				break
			case sig := <-signalChan:
				if sig.HasError() {
					return errors.Errorf("test_case %s: callback error: %s", tc.Name, sig.Error())
				}

				if sig.ImmuneCallBackID != uid {
					return errors.Errorf("test_case %s: incorrect callback_id: expected_callback_id '%s', got_callback_id '%s'", tc.Name, uid, sig.ImmuneCallBackID)
				}
				log.Infof("callback %d for test_case %s received", i, tc.Name)
			}
		}
	}

	return ex.dbTruncator.Truncate(ctx)
}

func (ex *Executor) sendRequest(ctx context.Context, r *request) (*response, error) {
	rb, err := json.Marshal(r.body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	req, err := http.NewRequestWithContext(ctx, r.method.String(), r.url, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", r.contentType)

	resp, err := ex.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &response{body: bytes.NewBuffer(buf), buf: buf, statusCode: resp.StatusCode}, nil
}
