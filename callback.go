package immune

import (
	"context"
	"fmt"
	"strings"
)

type CallbackConfiguration struct {
	MaxWaitSeconds uint   `json:"max_wait_seconds"`
	Port           uint   `json:"port"`
	Route          string `json:"route"`
	IDLocation     string `json:"id_location"`
}

const CallbackIDFieldName = "immune_callback_id"

// InjectCallbackID injects a callback id into field(expected to be a map[string]interface{} in r)
// in r, using CallbackIDFieldName as the key and value as the value.
func InjectCallbackID(field string, value interface{}, r M) error {
	// we may have separators referencing deeper fields in r e.g data.uid
	parts := strings.Split(field, ".")
	if len(parts) < 2 { // if it's less than 2, then there is no '.' in field
		v, ok := r[field]
		if !ok {
			return fmt.Errorf("the field %s, does not exist", field)
		}

		m, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("the field %s, is not an object in the request body", field)
		}
		m[CallbackIDFieldName] = value
	}

	nextLevel, err := getM(r, parts)
	if err != nil {
		return err
	}

	nextLevel[CallbackIDFieldName] = value

	return nil
}

type CallbackServer interface {
	ReceiveCallback() *Signal
	Start(ctx context.Context) error
	Stop()
}
