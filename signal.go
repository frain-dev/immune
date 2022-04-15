package immune

import "net/http"

// A Signal represents a single callback
type Signal struct {
	// ImmuneCallBackID collects the callback id from the request body, it's json tag
	// must always match immune.CallbackIDFieldName
	ImmuneCallBackID string `json:"immune_callback_id"`

	Request *http.Request // the http request that carried this callback signal
	Err     error
}

func (s *Signal) Error() string {
	return s.Err.Error()
}

func (s *Signal) HasError() bool {
	return s.Err != nil
}
