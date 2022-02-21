package immune

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type VariableMap struct {
	VariableToValue M
}

func (v VariableMap) GetString(key string) (string, bool) {
	value, ok := v.VariableToValue[key]
	if !ok {
		return "", false
	}

	str, ok := value.(string)
	if ok {
		return str, true
	}

	return fmt.Sprintf("%s", value), true
}

func (v VariableMap) Get(key string) (interface{}, bool) {
	value, ok := v.VariableToValue[key]
	return value, ok
}

func (v VariableMap) ProcessResponse(ctx context.Context, variableToField S, resp M) error {
	for varName, field := range variableToField {

		var value interface{}
		var ok bool

		// we may have separators referencing deeper fields in the response body e.g data.uid
		parts := strings.Split(field, ".")

		if len(parts) == 0 {
			value, ok = resp[field]
			if !ok {
				return errors.Errorf("variable %s's field %s not found in response", varName, field)
			}
		} else {
			next := M{}
			lastPart := parts[len(parts)-1]
			parts = parts[:len(parts)-1]

			track := ""
			for _, part := range parts {
				v := resp[part]

				next, ok = v.(map[string]interface{})
				if !ok {
					return errors.Errorf("the field %s, is not an object in response body", track)
				}
				track += part + "."
			}

			value = next[lastPart] // we have reached the last part of the "data.uid"
		}

		// only supporting string and int for now
		switch value.(type) {
		case string, int, int8, int32, int16, int64:
			break
		default:
			return errors.Errorf("variable %s is of type %T in the response body, only string & integer is currently supported", varName, value)
		}

		v.VariableToValue[varName] = value
	}

	return nil
}
