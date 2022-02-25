package immune

import (
	"strings"

	"github.com/pkg/errors"
)

type CallbackConfiguration struct {
	MaxWaitSeconds uint   `json:"max_wait_seconds"`
	Port           uint   `json:"port"`
	Route          string `json:"route"`
	IDLocation     string `json:"id_location"`
}

const CallbackIDFieldName = "immune_callback_id"

// InjectCallbackID injects a callback into field(expected to be a map[string]interface{} in r)
// in r, using CallbackIDFieldName as the key and v as the value.
func InjectCallbackID(field string, v interface{}, r M) error {
	// we may have separators referencing deeper fields in the response body e.g data.uid
	parts := strings.Split(field, ".")
	if len(parts) == 0 {
		v, ok := r[field]
		if !ok { // the field doesn't exist, so create it
			return errors.Errorf("the field %s, does not exist", field)
		}

		m, ok := v.(map[string]interface{})
		if !ok {
			return errors.Errorf("the field %s, is not an object in the request body", field)
		}
		m[CallbackIDFieldName] = v // we have reached the last part of the "data.uid"
	}

	nextLevel, err := getM(r, parts)
	if err != nil {
		return err
	}

	nextLevel[CallbackIDFieldName] = v // we have reached the last part of the "data.uid"

	return nil
}
