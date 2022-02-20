package net

import (
	"strings"

	"github.com/frain-dev/immune"
	"github.com/pkg/errors"
)

type request struct {
	body   immune.M
	url    string
	method immune.Method
}

func (r *request) processWithVariableMap(vm *immune.VariableMap) error {
	return r.traverse(r.body, vm)
}

func (r *request) traverse(m immune.M, vm *immune.VariableMap) error {
	for k, v := range m {
		switch v.(type) {
		case string:
			str := v.(string)
			if len(str) < 3 { // it must be in the format {x}
				continue
			}

			if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
				varName := str[1 : len(str)-2]
				value, exists := vm.Get(varName)
				if !exists {
					return errors.Errorf("variable %s does not exist in variable map")
				}

				m[k] = value // replace m[k] with the variable value
			}
		case immune.M: // TODO: may cause issues and have to change to  map[string]interface{}
			return r.traverse(v.(immune.M), vm)
		}
	}
	return nil
}
