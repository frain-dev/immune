package url

import (
	"strings"

	"github.com/frain-dev/immune"
	"github.com/pkg/errors"
)

type URL struct {
	variables []string
	url       string
}

// Parse parses an url string and records the segments containing variable references
// The returned URL object contains a slice of these segments and the original url.
//
// It works by looking for the first occurrence of the '{' character and then a corresponding '}'
// character repeatedly, recording all encountered variables until the end.
func Parse(s string) (*URL, error) {
	if len(s) == 0 {
		return nil, errors.New("url is empty")
	}

	u := &URL{variables: []string{}, url: s}

	// keep a set of encountered variables, since they will be all
	// replaced by the same value when ProcessWithVariableMap is called,
	// it makes sense to record just one instance if it
	set := map[string]bool{}
	for {
		v := nextVariable(s)
		if v == nil {
			break
		}
		s = s[v.closing:]

		if v.name == "" || set[v.name] { // if we've encountered this variable already or the segment is empty like this {}, continue
			continue
		}
		set[v.name] = true

		u.variables = append(u.variables, v.name)
	}

	return u, nil
}

// ProcessWithVariableMap replaces the url variable segments with their corresponding values from the variable map
func (u *URL) ProcessWithVariableMap(vm *immune.VariableMap) (string, error) {
	if len(u.variables) == 0 {
		return u.url, nil
	}

	result := u.url

	for _, variable := range u.variables {
		v, ok := vm.GetString(variable)
		if !ok {
			return "", errors.Errorf("variable %s not found in variable map", variable)
		}

		result = strings.Replace(result, "{"+variable+"}", v, -1)
	}

	return result, nil
}

type variableSegment struct {
	name    string
	closing int
}

// nextVariable looks for the first variable in the given string
func nextVariable(s string) *variableSegment {
	open := strings.IndexByte(s, '{')
	if open < 0 {
		return nil
	}

	closing := open
	for i, c := range s[open:] {
		if c == '}' {
			closing = open + i
			break
		}
	}

	if closing == open {
		panic("url: variable closing delimiter '}' is missing")
	}

	if s[open:closing+1] == "{}" {
		return &variableSegment{
			name:    "",
			closing: closing,
		}
	}

	return &variableSegment{
		name:    s[open+1 : closing],
		closing: closing,
	}
}
