package immune

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type VariableMap struct {
	variableToValue M
}

func (v VariableMap) GetString(key string) (string, bool) {
	value, ok := v.variableToValue[key]
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
	value, ok := v.variableToValue[key]
	return value, ok
}

func (v VariableMap) ProcessResponse(ctx context.Context, variableToField S, resp M) error {
	for varName, field := range variableToField {
		value, ok := resp[field]
		if !ok {
			return errors.Errorf("variable %s's field %s not found in response", varName, field)
		}

		// only supporting string and int for now
		switch value.(type) {
		case string, int, int8, int32, int16, int64:
			break
		default:
			return errors.Errorf("variable %s is of type %T in the response body, only string & integer is currently supported", varName, value)
		}

		v.variableToValue[varName] = value
	}

	return nil
}
