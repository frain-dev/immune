package immune

type CallbackConfiguration struct {
	MaxWaitSeconds uint   `json:"max_wait_seconds"`
	Port           uint   `json:"port"`
	Route          string `json:"route"`
	IDLocation     string `json:"id_location"`
}
