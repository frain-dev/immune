package fire

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Response struct {
	statusCode int
	buf        []byte
	body       *bytes.Buffer
}

type convoyResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (resp *Response) DecodeJSON(out interface{}) error {
	cr := convoyResponse{Data: out}
	err := json.NewDecoder(resp.body).Decode(&cr)
	if err != nil {
		return err
	}

	if !cr.Status {
		return fmt.Errorf("convoy error: %s, status code: %d", cr.Message, resp.statusCode)
	}

	return nil
}
