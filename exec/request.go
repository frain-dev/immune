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

			if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
				varName := value[1 : len(value)-1]
				value, exists := vm.Get(varName)
				if !exists {
					return errors.Errorf("variable %s does not exist in variable map", varName)
				}

				m[k] = value // replace m[k] with the variable value
			}
		case map[string]interface{}:
			// recursively traverse values with the type map[string]interface{}
			err := r.traverse(value, vm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
