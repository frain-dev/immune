package immune

import (
	"strings"

	"github.com/pkg/errors"
)

type Method string

var (
	DefaultSignatureHeader       = "X-Immune-Signature"
	DefaultEventIDHeader         = "X-Convoy-Event-ID"
	DefaultEventDeliveryIDHeader = "X-Convoy-EventDelivery-ID"
	DefaultHash                  = "SHA256"
	DefaultEncoding              = "hex"
)

const (
	MethodPost    Method = "POST"
	MethodPUT     Method = "PUT"
	MethodGet     Method = "GET"
	MethodPatch   Method = "PATCH"
	MethodHead    Method = "HEAD"
	MethodDelete  Method = "DELETE"
	MethodConnect Method = "CONNECT"
	MethodOptions Method = "OPTIONS"
	MethodTrace   Method = "TRACE"
)

func (m Method) IsValid() bool {
	switch m {
	case MethodPost,
		MethodPUT,
		MethodGet,
		MethodPatch,
		MethodHead,
		MethodDelete,
		MethodConnect,
		MethodOptions,
		MethodTrace:
		return true
	default:
		return false
	}
}

func (m Method) String() string {
	return string(m)
}

func (m *Method) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)

	*m = Method(str)

	if !m.IsValid() {
		return errors.Errorf("unknown http method %s", m.String())
	}

	return nil
}
