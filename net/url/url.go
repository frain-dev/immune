package url

import (
	"strings"

	"github.com/frain-dev/immune"
	"github.com/pkg/errors"
)

type URL struct {
	segments []*segment
	url      string
}

func Parse(s string) (*URL, error) {
	if len(s) == 0 {
		return nil, errors.New("url is empty")
	}

	u := &URL{segments: []*segment{}, url: s}

	for {
		seg := nextSegment(s)
		if seg == nil {
			break
		}

		u.segments = append(u.segments, seg)
		s = s[seg.end:]
	}

	return u, nil
}

func (u *URL) ProcessWithVariableMap(vm *immune.VariableMap) (string, error) {
	if len(u.segments) == 0 {
		return u.url, nil
	}

	var result string

	url := u.url

	for _, s := range u.segments {
		v, ok := vm.GetString(s.name)
		if !ok {
			return "", errors.Errorf("variable %s not found in variable map", s.name)
		}

		result = strings.Replace(url, url[s.start:s.end+1], v, -1)
		url = url[s.end+1:]
	}

	return result, nil
}

type segment struct {
	start int
	end   int
	name  string
}

func nextSegment(s string) *segment {
	open := strings.IndexRune(s, '{')
	if open < 0 {
		return nil
	}

	close := open
	for i, c := range s[open:] {
		if c == '}' {
			close = open + i
			break
		}
	}

	if close == open {
		panic("url: variable closing delimiter '}' is missing")
	}

	return &segment{
		start: open,
		end:   close,
		name:  s[open+1 : close],
	}
}
