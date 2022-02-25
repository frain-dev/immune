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

// GetString gets the value of key from the variable map, if the value
// isn't of the string type, it will be converted to string via fmt.Sprintf
// and returned
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

// Get gets the value of key from the variable map
func (v VariableMap) Get(key string) (interface{}, bool) {
	value, ok := v.VariableToValue[key]
	return value, ok
}

// ProcessResponse takes the variables declared in variableToField from values, and stores them in the
// variable map.
func (v VariableMap) ProcessResponse(ctx context.Context, variableToField S, values M) error {
	for varName, field := range variableToField {

		value, err := getKeyInMap(field, values)
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

// getM fetches the item in parts, going one level deeper with each iteration of parts
// the result of each iteration is expected to be of type map[string]interface{}
func getM(m M, parts []string) (M, error) {
	nextLevel := m
	var ok bool
	var v interface{}

	track := ""
	for _, part := range parts {
		v, ok = nextLevel[part]
		if !ok {
			return nil, errors.Errorf("the field %s, does not exist", track+part) // avoid printing the trailing dot
		}

		nextLevel, ok = v.(map[string]interface{})
		if !ok {
			return nil, errors.Errorf("the field %s, is not an object in response body", track+part)
		}

		track += part + "."
	}

	return nextLevel, nil
}
