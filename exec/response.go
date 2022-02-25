package exec

import (
	"bytes"
	"encoding/json"
)

type response struct {
	statusCode int
	body       *bytes.Buffer
}

func (resp *response) Decode(out interface{}) error {
	return json.NewDecoder(resp.body).Decode(out)
}
