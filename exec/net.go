package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/url"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Executor is used to execute tests
type Executor struct {
	callbackIDLocation     string
	baseURL                string
	maxCallbackWaitSeconds uint
	s                      immune.CallbackServer
	client                 *http.Client
	vm                     *immune.VariableMap
}

func NewExecutor(s immune.CallbackServer, client *http.Client, vm *immune.VariableMap, maxCallbackWaitSeconds uint, baseURL string, callbackIDLocation string) *Executor {
	return &Executor{
		s:                      s,
		vm:                     vm,
		client:                 client,
		baseURL:                baseURL,
		callbackIDLocation:     callbackIDLocation,
		maxCallbackWaitSeconds: maxCallbackWaitSeconds,
	}
}

// ExecuteSetupTestCase executes setup test cases
func (ex *Executor) ExecuteSetupTestCase(ctx context.Context, setupTC *immune.SetupTestCase) error {
	u, err := url.Parse(fmt.Sprintf("%s%s", ex.baseURL, setupTC.Endpoint))
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %d: failed to parse url", setupTC.Position)
	}

	result, err := u.ProcessWithVariableMap(ex.vm)
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %d: failed to process parsed url with variable map", setupTC.Position)
	}

	r := &request{
		contentType: "application/json",
		url:         result,
		body:        setupTC.RequestBody,
		method:      setupTC.HTTPMethod,
	}

	err = r.processWithVariableMap(ex.vm)
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %d: failed to process request body with variable map", setupTC.Position)
	}
	//log.Infof("setup_test_case %d: request body: %s", setupTC.Position, pretty.Sprint(r.body))

	resp, err := ex.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	if setupTC.StatusCode != resp.statusCode {
		return errors.Errorf("setup_test_case %d: wants status code %d but got status code %d", setupTC.Position, setupTC.StatusCode, resp.statusCode)
	}

	if setupTC.ResponseBody {
		if resp.body.Len() == 0 {
			return errors.Errorf("setup_test_case %d: wants response body but got no response body", setupTC.Position)
		}

		m := immune.M{}
		err = resp.Decode(&m)
		if err != nil {
			return errors.Wrapf(err, "setup_test_case %d: failed to decode response body", setupTC.Position)
		}

		if setupTC.StoreResponseVariables != nil {
			err = ex.vm.ProcessResponse(ctx, setupTC.StoreResponseVariables, m)
			if err != nil {
				return errors.Wrapf(err, "setup_test_case %d: failed to process response body", setupTC.Position)
			}
		}
	} else {
		if resp.body.Len() > 0 {
			return errors.Errorf("setup_test_case %d: does not want a response body but got a response body: '%s'", setupTC.Position, resp.body.String())
		}
	}

	return nil
}

// ExecuteTestCase executes test cases, it waits for callback if necessary
func (ex *Executor) ExecuteTestCase(ctx context.Context, tc *immune.TestCase) error {
	u, err := url.Parse(fmt.Sprintf("%s%s", ex.baseURL, tc.Endpoint))
	if err != nil {
		return errors.Wrapf(err, "test_case %d: failed to parse url", tc.Position)
	}

	result, err := u.ProcessWithVariableMap(ex.vm)
	if err != nil {
		return errors.Wrapf(err, "test_case %d: failed to process parsed url with variable map", tc.Position)
	}

	var uid string
	if tc.Callback.Enabled {
		uid = uuid.New().String()
		err = immune.InjectCallbackID(ex.callbackIDLocation, uid, tc.RequestBody)
		if err != nil {
			return errors.Wrapf(err, "test_case %d: failed to inject callback id into request body", tc.Position)
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
		return errors.Wrapf(err, "test_case %d: failed to process request body with variable map", tc.Position)
	}

	resp, err := ex.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	if tc.StatusCode != resp.statusCode {
		return errors.Errorf("test_case %d: wants status code %d but got status code %d", tc.Position, tc.StatusCode, resp.statusCode)
	}

	if tc.ResponseBody {
		if resp.body.Len() == 0 {
			return errors.Errorf("test_case %d: wants response body but got no response body: %+v", tc.Position, resp)
		}

		m := immune.M{}
		err = resp.Decode(&m)
		if err != nil {
			return errors.Errorf("test_case %d: failed to decode response body: %+v", tc.Position, resp)
		}

	} else {
		if resp.body.Len() > 0 {
			return errors.Wrapf(err, "test_case %d: does not want a response body but got a response body: '%s'", tc.Position, resp.body.String())
		}
	}

	if tc.Callback.Enabled {
		cctx, cancel := context.WithTimeout(context.Background(), time.Duration(ex.maxCallbackWaitSeconds)*time.Second)
		defer cancel()

		for i := uint(1); i <= tc.Callback.Times; i++ {
			select {
			case <-cctx.Done():
				log.Infof("succesfully received %d callbacks for test_case %d before max callback wait seconds elapsed", i, tc.Position)
				break
			default:
				sig := ex.s.ReceiveCallback()
				if sig.ImmuneCallBackID != uid {
					return errors.Errorf("test_case %d: incorrect callback_id: expected_callback_id '%s', got_callback_id '%s'", tc.Position, uid, sig.ImmuneCallBackID)
				}
				log.Infof("callback %d for test_case %d received", i, tc.Position)
			}
		}
	}

	return nil
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

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &response{body: bytes.NewBuffer(buf), statusCode: resp.StatusCode}, nil
}
