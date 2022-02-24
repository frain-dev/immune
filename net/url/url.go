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

// Parse parses an url string and records the segments containing variable references
// The returned URL object contains a slice of these segments and the original url.
//
// It works by looking for the first occurrence of the '{' character and then a corresponding '}'
// character, when the first segment is found the ending index is added to i, the same is done for
// subsequent segments. By pushing i forward as we iterate,we can track how far along in the original
// string we have come, when iterating for multiple segments.
func Parse(s string) (*URL, error) {
	if len(s) == 0 {
		return nil, errors.New("url is empty")
	}

	u := &URL{segments: []*segment{}, url: s}
	i := 0
	for {
		seg := nextSegment(s)
		if seg == nil {
			break
		}

		s = s[seg.end:] // cut s so the next call to nextSegment has a smaller string to run over

		// add the last recorded index addition
		seg.end = seg.end + i
		seg.start = seg.start + i
		i += seg.end // increment the index again

		u.segments = append(u.segments, seg)
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

	return &segment{
		start: open,
		end:   closing,
		name:  s[open+1 : closing],
	}
}
