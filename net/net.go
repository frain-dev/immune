package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/callback"
	"github.com/frain-dev/immune/net/url"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Net struct {
	s                      *callback.Server
	client                 *http.Client
	vm                     *immune.VariableMap
	maxCallbackWaitSeconds int
	baseURL                string
}

func (n *Net) ExecuteSetupTestCase(ctx context.Context, setupTC *immune.SetupTestCase) error {
	u, err := url.Parse(fmt.Sprintf("%s%s", n.baseURL, setupTC.Endpoint))
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %d: failed to parse url", setupTC.Position)
	}

	result, err := u.ProcessWithVariableMap(n.vm)
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %d: failed to process parsed url with variable map", setupTC.Position)
	}

	r := &request{
		url:    result,
		body:   setupTC.RequestBody,
		method: setupTC.HTTPMethod,
	}

	err = r.processWithVariableMap(n.vm)
	if err != nil {
		return errors.Wrapf(err, "setup_test_case %d: failed to process request body with variable map", setupTC.Position)
	}

	resp, err := n.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	if setupTC.ResponseBody {
		if len(resp.body) == 0 {
			return errors.Wrapf(err, "setup_test_case %d: wants response body but got no response body", setupTC.Position)
		}

		m := immune.M{}
		err = resp.Decode(&m)
		if err != nil {
			return err
		}

		if setupTC.StoreResponseVariables != nil {
			err = n.vm.ProcessResponse(ctx, setupTC.StoreResponseVariables, m)
			if err != nil {
				return err
			}
		}
	} else {
		if len(resp.body) > 0 {
			return errors.Wrapf(err, "setup_test_case %d: does not want a response body but got a response body: '%s'", setupTC.Position, string(resp.body))
		}
	}

	return nil
}

func (n *Net) ExecuteTestCase(ctx context.Context, tc *immune.TestCase) error {
	u, err := url.Parse(fmt.Sprintf("%s%s", n.baseURL, tc.Endpoint))
	if err != nil {
		return errors.Wrapf(err, "test_case %d: failed to parse url", tc.Position)
	}

	result, err := u.ProcessWithVariableMap(n.vm)
	if err != nil {
		return errors.Wrapf(err, "test_case %d: failed to process parsed url with variable map", tc.Position)
	}

	r := &request{
		body:   tc.RequestBody,
		url:    result,
		method: tc.HTTPMethod,
	}

	err = r.processWithVariableMap(n.vm)
	if err != nil {
		return errors.Wrapf(err, "test_case %d: failed to process request body with variable map", tc.Position)
	}

	resp, err := n.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	if tc.ResponseBody {
		if len(resp.body) == 0 {
			return errors.Wrapf(err, "test_case %d: wants response body but got no response body", tc.Position)
		}

		m := immune.M{}
		err = resp.Decode(&m)
		if err != nil {
			return err
		}
	} else {
		if len(resp.body) > 0 {
			return errors.Wrapf(err, "test_case %d: does not want a response body but got a response body: '%s'", tc.Position, string(resp.body))
		}
	}

	if tc.Callback.Enabled {
		cctx, cancel := context.WithTimeout(context.Background(), time.Duration(n.maxCallbackWaitSeconds)*time.Second)
		defer cancel()

		for i := 1; i <= tc.Callback.Times; i++ {
			select {
			case <-cctx.Done():
				log.Infof("succesfully received %d callbacks for test_case %d before max callback wait seconds elapsed", i, tc.Position)
				break
			default:
				n.s.ReceiveCallback()
				log.Infof("callback %d for test_case %d received", i, tc.Position)
			}
		}
	}

	return nil
}

func (n *Net) sendRequest(ctx context.Context, r *request) (*response, error) {
	rb, err := json.Marshal(r.body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	req, err := http.NewRequestWithContext(ctx, r.method.String(), r.url, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf []byte
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &response{body: buf, statusCode: resp.StatusCode}, nil
}
