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

		value, err := getKeyInMap(field, resp)
		if err != nil {
			return err
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

func getKeyInMap(field string, resp M) (interface{}, error) {
	var value interface{}
	var ok bool

	// we may have separators referencing deeper fields in the response body e.g data.uid
	parts := strings.Split(field, ".")

	if len(parts) == 0 {
		value, ok = resp[field]
		if !ok {
			return nil, errors.Errorf("field %s not found in response", field)
		}
	} else {
		lastPart := parts[len(parts)-1]
		parts = parts[:len(parts)-1]

		nextLevel, err := getM(resp, parts)
		if err != nil {
			return nil, err
		}

		value = nextLevel[lastPart] // we have reached the last part of the "data.uid"
	}

	return value, nil
}

const CallbackIDFieldName = "immune_callback_id"

func InjectCallbackID(field string, value interface{}, resp M) error {
	// we may have separators referencing deeper fields in the response body e.g data.uid
	parts := strings.Split(field, ".")
	if len(parts) == 0 {
		v, ok := resp[field]
		if !ok { // the field doesn't exist, so create it
			v = map[string]interface{}{}
			resp[field] = v
		}

		m, ok := v.(map[string]interface{})
		if !ok {
			return errors.Errorf("the field %s, is not an object in the request body", field)
		}
		m[CallbackIDFieldName] = value // we have reached the last part of the "data.uid"
	}

	nextLevel, err := getM(resp, parts)
	if err != nil {
		return err
	}

	nextLevel[CallbackIDFieldName] = value // we have reached the last part of the "data.uid"

	return nil
}

func getM(m M, parts []string) (M, error) {
	nextLevel := M{}
	var ok bool

	track := ""
	for _, part := range parts {
		nextLevel, ok = m[part].(map[string]interface{})
		if !ok {
			return nil, errors.Errorf("the field %s, is not an object in response body", track[:len(track)-1]) // avoid printing the trailing dot
		}
		track += part + "."
	}

	return nextLevel, nil
}
