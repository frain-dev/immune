package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/frain-dev/immune/callback"

	"github.com/frain-dev/immune"
)

type Net struct {
	client  *http.Client
	s       *callback.Server
	baseURL string
}

func (n *Net) ExecuteSetupTestCase(ctx context.Context, setupTC *immune.SetupTestCase, vm immune.VariableMap) error {
	r := &request{
		body:   setupTC.RequestBody,
		url:    fmt.Sprintf("%s%s", n.baseURL, setupTC.Endpoint),
		method: setupTC.HTTPMethod,
	}

	resp, err := n.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	report := &immune.SetupTestCaseReport{
		WantsResponseBody: setupTC.ResponseBody,
	}

	if setupTC.ResponseBody {
		if len(resp.body) == 0 {
			report.HasResponseBody = false
			setupTC.Report = report
		}

		m := immune.M{}
		err = resp.Decode(&m)
		if err != nil {
			return err
		}

		err = vm.ProcessResponse(ctx, setupTC.StoreResponseVariables, m)
		if err != nil {
			return err
		}
	} else {
		if len(resp.body) > 0 {
			report.HasResponseBody = true
			setupTC.Report = report
		}
	}

	return nil
}

func (n *Net) ExecuteTestCase(ctx context.Context, tc *immune.TestCase) {

}

type request struct {
	body   immune.M
	url    string
	method immune.Method
}

type response struct {
	body []byte
}

func (resp *response) Decode(out interface{}) error {
	return json.Unmarshal(resp.body, out)
}

func (n *Net) sendRequest(ctx context.Context, r *request) (*response, error) {
	rb, err := json.Marshal(r.body)
	req, err := http.NewRequestWithContext(ctx, r.method.String(), r.url, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("status", resp.StatusCode).Infof("non ok status code returned")
	}

	var buf []byte
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &response{body: buf}, nil
}
