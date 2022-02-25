package immune

type SetupTestCase struct {
	StoreResponseVariables S      `json:"store_response_variables"`
	RequestBody            M      `json:"request_body"`
	ResponseBody           bool   `json:"response_body"`
	Endpoint               string `json:"endpoint"`
	HTTPMethod             Method `json:"http_method"`
	Position               int    `json:"-"`
	//Report                 *SetupTestCaseReport `json:"-"`
}

type SetupTestCaseReport struct {
	WantsResponseBody bool
	HasResponseBody   bool
}

type TestCase struct {
	Position     int      `json:"-"`
	HTTPMethod   Method   `json:"http_method"`
	Endpoint     string   `json:"endpoint"`
	ResponseBody bool     `json:"response_body"`
	Callback     Callback `json:"callback"`
	RequestBody  M        `json:"request_body"`
}

type Callback struct {
	Enabled bool `json:"enabled"`
	Times   uint `json:"times"`
}
