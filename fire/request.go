package fire

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/frain-dev/immune"
)

type Request struct {
	contentType string
	url         string
	method      immune.Method
	body        []byte
	headers     http.Header
}

func (r *Request) WithJSONBody(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.body = b
	return nil
}

func (r *Request) AddHeader(key, val string) {
	if r.headers == nil {
		r.headers = http.Header{}
	}

	r.headers.Add(key, val)
}

func (r *Request) SendRequest(ctx context.Context) (*Response, error) {
	bb := bytes.NewBuffer(r.body)

	req, err := http.NewRequestWithContext(ctx, r.method.String(), r.url, bb)
	if err != nil {
		return nil, err
	}

	for k, v := range r.headers {
		req.Header[k] = v
	}

	req.Header.Add("Content-Type", r.contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{body: bytes.NewBuffer(buf), statusCode: resp.StatusCode}, nil
}

//// processWithVariableMap replaces all variable references in the request body with
//// their corresponding values from the variable map
//func (r *Request) processWithVariableMap(vm *immune.VariableMap) error {
//	return r.traverse(r.body, vm)
//}

// traverse examines all key-value pairs in m replacing all values that reference
// variables with their corresponding values from the variable map
//func (r *Request) traverse(m immune.M, vm *immune.VariableMap) error {
//	for k, v := range m {
//		switch value := v.(type) {
//		case string: // only string values in the map can reference variables in the format "{variable_name}"
//			if len(value) < 3 { // at least three characters must be present in the format {x}
//				continue
//			}
//
//			val, err := getVariableValue(value, vm)
//			if err != nil {
//				return err
//			}
//
//			m[k] = val // replace m[k] with the variable value
//
//		case map[string]interface{}:
//			// recursively traverse values with the type map[string]interface{}
//			err := r.traverse(value, vm)
//			if err != nil {
//				return err
//			}
//		case []interface{}:
//			for i, sliceValue := range value {
//				var val interface{}
//				var err error
//
//				switch s := sliceValue.(type) {
//				case string:
//					val, err = getVariableValue(s, vm)
//					if err != nil {
//						return err
//					}
//					value[i] = val
//				case map[string]interface{}:
//					// recursively traverse values with the type map[string]interface{}
//					err := r.traverse(s, vm)
//					if err != nil {
//						return err
//					}
//					continue
//				}
//			}
//		}
//	}
//	return nil
//}

//func getVariableValue(str string, vm *immune.VariableMap) (interface{}, error) {
//	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
//		varName := str[1 : len(str)-1]
//		val, exists := vm.Get(varName)
//		if !exists {
//			return nil, errors.Errorf("variable %s does not exist in variable map", varName)
//		}
//		return val, nil
//	}
//
//	return str, nil // return original since, it's not a variable reference
//}
