package exec

import (
	"strings"

	"github.com/frain-dev/immune"
	"github.com/pkg/errors"
)

type request struct {
	contentType string
	url         string
	method      immune.Method
	body        immune.M
}

// processWithVariableMap replaces all variable references in the request body with
// their corresponding values from the variable map
func (r *request) processWithVariableMap(vm *immune.VariableMap) error {
	return r.traverse(r.body, vm)
}

// traverse examines all key-value pairs in m replacing all values that reference
// variables with their corresponding values from the variable map
func (r *request) traverse(m immune.M, vm *immune.VariableMap) error {
	for k, v := range m {
		switch value := v.(type) {
		case string: // only string values in the map can reference variables in the format "{variable_name}"
			if len(value) < 3 { // at least three characters must be present in the format {x}
				continue
			}

			val, err := getVariableValue(value, vm)
			if err != nil {
				return err
			}

			m[k] = val // replace m[k] with the variable value

		case map[string]interface{}:
			// recursively traverse values with the type map[string]interface{}
			err := r.traverse(value, vm)
			if err != nil {
				return err
			}
		case []interface{}:
			for i, sliceValue := range value {
				var val interface{}
				var err error

				switch s := sliceValue.(type) {
				case string:
					val, err = getVariableValue(s, vm)
					if err != nil {
						return err
					}
					value[i] = val
				case map[string]interface{}:
					// recursively traverse values with the type map[string]interface{}
					err := r.traverse(s, vm)
					if err != nil {
						return err
					}
					continue
				}
			}
		}
	}
	return nil
}

func getVariableValue(str string, vm *immune.VariableMap) (interface{}, error) {
	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		varName := str[1 : len(str)-1]
		val, exists := vm.Get(varName)
		if !exists {
			return nil, errors.Errorf("variable %s does not exist in variable map", varName)
		}
		return val, nil
	}

	return str, nil // return original since, it's not a variable reference
}
