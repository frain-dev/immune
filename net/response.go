package net

import "encoding/json"

type response struct {
	statusCode int
	body       []byte
}

func (resp *response) Decode(out interface{}) error {
	return json.Unmarshal(resp.body, out)
}
