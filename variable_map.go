package immune

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type VariableMap struct {
	VariableToValue M
}

func NewVariableMap() *VariableMap {
	return &VariableMap{VariableToValue: M{}}
}

// GetString gets the value of key from the variable map, if the value
// isn't of the string type, it will be converted to string via fmt.Sprintf
// and returned
func (v *VariableMap) GetString(key string) (string, bool) {
	value, ok := v.VariableToValue[key]
	if !ok {
		return "", false
	}

	str, ok := value.(string)
	if ok {
		return str, true
	}

	return fmt.Sprintf("%v", value), true
}

// Get gets the value of key from the variable map
func (v *VariableMap) Get(key string) (interface{}, bool) {
	value, ok := v.VariableToValue[key]
	return value, ok
}

// ProcessResponse takes the variables declared in variableToField from values, and stores them in the
// variable map.
func (v *VariableMap) ProcessResponse(ctx context.Context, variableToField S, values M) error {
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
			return fmt.Errorf("variable %s is of type %T in the response body, only string & integer is currently supported", varName, value)
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
	if len(parts) < 2 { // if it's less than 2, then there is no '.' in field
		if isArray(field) {
			var err error
			value, err = getArrayValue(field, resp)
			if err != nil {
				return nil, fmt.Errorf("field %s: %v", field, err)
			}
		} else {
			value, ok = resp[field]
			if !ok {
				return nil, fmt.Errorf("field %s: not found", field)
			}
		}
	} else {
		lastPart := parts[len(parts)-1]
		parts = parts[:len(parts)-1]

		nextLevel, err := getM(resp, parts)
		if err != nil {
			return nil, err
		}

		value, ok = nextLevel[lastPart] // we have reached the last part of the "data.uid"
		if !ok {
			return nil, fmt.Errorf("field %s: not found", field)
		}
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
		if isArray(part) {
			var err error
			v, err = getArrayValue(part, nextLevel)
			if err != nil {
				return nil, fmt.Errorf("field %s: %v", track+part, err)
			}
		} else {
			v, ok = nextLevel[part]
			if !ok {
				return nil, fmt.Errorf("field %s: not found", track+part)
			}
		}

		nextLevel, ok = v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field %s: required type is object but got %T", track+part, v)
		}

		track += part + "."
	}

	return nextLevel, nil
}

func isArray(v string) bool {
	return strings.Contains(v, "[") && strings.HasSuffix(v, "]")
}

// getArrayValue parses a variable reference and returns the value in m
// v is expected to be of format field[index] e.g. tunnels[0]
func getArrayValue(v string, m M) (interface{}, error) {
	open := strings.Index(v, "[")
	closer := strings.Index(v, "]")

	ixStr := v[open+1 : closer]
	ix, err := strconv.Atoi(ixStr)
	if err != nil {
		return nil, fmt.Errorf("invalid index notation: %s", ixStr)
	}

	if ix < 0 {
		return nil, fmt.Errorf("invalid index range: %d", ix)
	}

	name := v[:open]
	field, ok := m[name]
	if !ok {
		return nil, errors.New("not found")
	}

	sliceValue, ok := field.([]interface{})
	if !ok {
		return nil, fmt.Errorf("required type is an array but has type %T", field)
	}

	if len(sliceValue) < ix+1 {
		return nil, fmt.Errorf("index out of range with length %d", len(sliceValue))
	}

	return sliceValue[ix], nil
}
