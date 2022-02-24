package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/frain-dev/immune"
	log "github.com/sirupsen/logrus"
)

type System struct {
	BaseURL        string                 `json:"base_url"`
	Callback       CallbackConfiguration  `json:"callback"`
	Variables      *immune.VariableMap    `json:"-"`
	SetupTestCases []immune.SetupTestCase `json:"setup_test_cases"`
	TestCases      []immune.TestCase      `json:"test_cases"`
	needsCallback  bool
}

type CallbackConfiguration struct {
	MaxWaitSeconds uint   `json:"max_wait_seconds"`
	Port           uint   `json:"port"`
	Route          string `json:"route"`
	IDLocation     string `json:"id_location"`
}

func NewSystem(filePath string) (*System, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	sys := &System{}

	err = json.NewDecoder(f).Decode(sys)
	if err != nil {
		return nil, err
	}

	sys.Variables = &immune.VariableMap{VariableToValue: immune.M{}}
	return sys, nil
}

const maxCallbackWait = 5

func (s *System) Clean() error {
	if s.BaseURL == "" {
		return errors.New("base url cannot be empty")
	}

	_, err := url.Parse(s.BaseURL)
	if err != nil {
		return fmt.Errorf("base url is not a vaild url: %v", err)
	}

	if s.Callback.MaxWaitSeconds == 0 {
		log.Warnf("max callback wait seconds is 0, using default value of %d seconds", maxCallbackWait)
		s.Callback.MaxWaitSeconds = maxCallbackWait
	}

	varNameToSetupTC := map[string]int{}

	for i := range s.SetupTestCases {
		tc := &s.SetupTestCases[i]
		tcNum := i + 1

		// ensure no variable name is used twice
		for name := range tc.StoreResponseVariables {
			ix, ok := varNameToSetupTC[name]
			if ok {
				return fmt.Errorf("setup_test_case %d: variable name %s already used in setup_test_case %d", tcNum, name, ix)
			}

			varNameToSetupTC[name] = tcNum
		}

		if len(tc.Endpoint) == 0 {
			return fmt.Errorf("setup_test_case %d: endpoint cannot be empty", tcNum)
		}

		if !strings.HasPrefix(tc.Endpoint, "/") {
			return fmt.Errorf("setup_test_case %d: endpoint must begin with /", tcNum)
		}

		if !tc.HTTPMethod.IsValid() {
			return fmt.Errorf("setup_test_case %d: invalid method: %s", tcNum, tc.HTTPMethod.String())
		}

		s.SetupTestCases[i].Position = tcNum
	}

	for i := range s.TestCases {
		tc := &s.TestCases[i]
		tcNum := i + 1

		if len(tc.Endpoint) == 0 {
			return fmt.Errorf("test_case %d: endpoint cannot be empty", tcNum)
		}

		if !strings.HasPrefix(tc.Endpoint, "/") {
			return fmt.Errorf("test_case %d: endpoint must begin with /", tcNum)
		}

		if !tc.HTTPMethod.IsValid() {
			return fmt.Errorf("test_case %d: invalid method: %s", tcNum, tc.HTTPMethod.String())
		}

		if tc.Callback.Enabled {
			s.needsCallback = true

			if tc.Callback.Times == 0 {
				return fmt.Errorf("test_case %d: if callback is enabled then times must be greater than 0", tcNum)
			}
		}

		s.TestCases[i].Position = tcNum
	}

	return nil
}

func (s *System) NeedsCallbackServer() bool {
	return s.needsCallback
}
