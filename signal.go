package immune

// A Signal represents a single callback
type Signal struct {
	// ImmuneCallBackID collects the callback id from the request body, it's json tag
	// must always match immune.CallbackIDFieldName
	ImmuneCallBackID string `json:"immune_callback_id"`
}
