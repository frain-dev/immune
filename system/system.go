package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/frain-dev/immune"
)

type System struct {
	BaseURL                string `json:"base_url"`
	needsCallback          bool
	MaxCallbackWaitSeconds uint                   `json:"max_callback_wait_seconds"`
	Variables              *immune.VariableMap    `json:"-"`
	SetupTestCases         []immune.SetupTestCase `json:"setup_test_cases"`
	TestCases              []immune.TestCase      `json:"test_cases"`
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

const maxCallbackWait = 2

func (s *System) Clean() error {
	if s.BaseURL == "" {
		return errors.New("base url cannot be empty")
	}

	_, err := url.Parse(s.BaseURL)
	if err != nil {
		return fmt.Errorf("base url is not a vaild url: %v", err)
	}

	if s.MaxCallbackWaitSeconds == 0 {
		log.Warnf("max callback wait seconds is 0, using default value of %d seconds", maxCallbackWait)
		s.MaxCallbackWaitSeconds = maxCallbackWait
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
