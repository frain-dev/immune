package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/frain-dev/immune"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// System represents the entire suite to be run against an API
type System struct {
	BaseURL        string                       `json:"base_url"`
	EventTargetURL string                       `json:"event_target_url" envconfig:"IMMUNE_EVENT_TARGET_URL"`
	Database       immune.Database              `json:"database"`
	Callback       immune.CallbackConfiguration `json:"callback"`
	Variables      *immune.VariableMap          `json:"-"`
	SetupTestCases []immune.SetupTestCase       `json:"setup_test_cases"`
	TestCases      []immune.TestCase            `json:"test_cases"`
	needsCallback  bool
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

	envOverride := &System{}
	err = envconfig.Process("IMMUNE", envOverride)
	if err != nil {
		return nil, err
	}

	processOverride(sys, envOverride)

	sys.Variables = &immune.VariableMap{VariableToValue: immune.M{}}
	return sys, nil
}

func processOverride(sys, override *System) {
	if override.EventTargetURL != "" {
		sys.EventTargetURL = override.EventTargetURL
	}

	if _, ok := os.LookupEnv("IMMUNE_SSL"); ok {
		sys.Callback.SSL = override.Callback.SSL
	}

	if override.Callback.SSLKeyFile != "" {
		sys.Callback.SSLKeyFile = override.Callback.SSLKeyFile
	}

	if override.Callback.SSLCertFile != "" {
		sys.Callback.SSLCertFile = override.Callback.SSLCertFile
	}
}

const maxCallbackWait = 5

// Clean validates the System's data
func (s *System) Clean() error {
	if s.BaseURL == "" {
		return errors.New("base url cannot be empty")
	}

	if s.Callback.SSL {
		if s.Callback.SSLCertFile == "" || s.Callback.SSLKeyFile == "" {
			return errors.New("both cert_file and key_file are required for ssl")
		}
	}

	_, err := url.Parse(s.BaseURL)
	if err != nil {
		return fmt.Errorf("base url is not a vaild url: %v", err)
	}

	if s.Callback.MaxWaitSeconds == 0 {
		log.Warnf("max callback wait seconds is 0, using default value of %d seconds", maxCallbackWait)
		s.Callback.MaxWaitSeconds = maxCallbackWait
	}

	for i := range s.TestCases {
		tc := &s.TestCases[i]

		if tc.Name == "" {
			return errors.New("test case name cannot be empty")
		}

		if len(tc.Endpoint) == 0 {
			return fmt.Errorf("test_case %s: endpoint cannot be empty", tc.Name)
		}

		if tc.StatusCode < 100 || tc.StatusCode > 511 {
			return fmt.Errorf("test_case %s: valid range for status_code is 100-511", tc.Name)
		}

		if !strings.HasPrefix(tc.Endpoint, "/") {
			return fmt.Errorf("test_case %s: endpoint must begin with /", tc.Name)
		}

		if !tc.HTTPMethod.IsValid() {
			return fmt.Errorf("test_case %s: invalid method: %s", tc.Name, tc.HTTPMethod.String())
		}

		if tc.Callback.Enabled {
			s.needsCallback = true

			if tc.Callback.Times == 0 {
				return fmt.Errorf("test_case %s: if callback is enabled then times must be greater than 0", tc.Name)
			}
		}
	}

	return nil
}

func (s *System) NeedsCallbackServer() bool {
	return s.needsCallback
}
