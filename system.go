package immune

type System struct {
	BaseURL                string          `json:"base_url"`
	MaxCallbackWaitSeconds int             `json:"max_callback_wait_seconds"`
	Variables              VariableMap     `json:"-"`
	SetupTestCases         []SetupTestCase `json:"setup_test_cases"`
	TestCases              []TestCase      `json:"test_cases"`
}

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
	Enabled bool
	Times   int
}
