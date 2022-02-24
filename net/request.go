package net

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
				varName := str[1 : len(str)-1]
				value, exists := vm.Get(varName)
				if !exists {
					return errors.Errorf("variable %s does not exist in variable map", varName)
				}

				m[k] = value // replace m[k] with the variable value
			}
		case map[string]interface{}: // TODO: may cause issues and have to change to  map[string]interface{}
			return r.traverse(v.(immune.M), vm)
		}
	}
	return nil
}
