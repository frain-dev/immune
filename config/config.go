package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/frain-dev/convoy/util"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type TestType string

const (
	IngestTest  TestType = "ingest"
	FanOutTest  TestType = "fan_out"
	SingleTest  TestType = "single"
	DynamicTest TestType = "dynamic"
	PubSubTest  TestType = "pub_sub"
)

// Config represents the entire suite to be run against an API
type Config struct {
	ConvoyURL      string   `json:"convoy_url" envconfig:"IMMUNE_CONVOY_URL"`
	EndpointURL    string   `json:"endpoint_url" envconfig:"IMMUNE_ENDPOINT_URL"`
	EndpointSecret string   `json:"endpoint_secret" envconfig:"IMMUNE_ENDPOINT_SECRET"`
	Events         int64    `json:"events" envconfig:"IMMUNE_EVENTS"`
	TestType       TestType `json:"test_type" envconfig:"IMMUNE_CONVOY_URL"`
	LogFile        string   `json:"log_file" envconfig:"IMMUNE_LOG_FILE"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	envOverride := &Config{}
	err = envconfig.Process("IMMUNE", envOverride)
	if err != nil {
		return nil, err
	}

	Override(cfg, envOverride)

	return cfg, nil
}

func (t TestType) IsValid() bool {
	switch t {
	case IngestTest,
		FanOutTest,
		SingleTest,
		DynamicTest,
		PubSubTest:
		return true
	default:
		return false
	}
}

//func processOverride(sys, override *Config) {
//	if override.EventTargetURL != "" {
//		sys.EventTargetURL = override.EventTargetURL
//	}
//
//	if _, ok := os.LookupEnv("IMMUNE_SSL"); ok {
//		sys.Callback.SSL = override.Callback.SSL
//	}
//
//	if override.Callback.SSLKeyFile != "" {
//		sys.Callback.SSLKeyFile = override.Callback.SSLKeyFile
//	}
//
//	if override.Callback.SSLCertFile != "" {
//		sys.Callback.SSLCertFile = override.Callback.SSLCertFile
//	}
//}

const numEvents = 5000

// Validate validates the Config's data
func (c *Config) Validate() error {
	if c.Events == 0 {
		log.Warnf("number of events to send/expect is zero, using default value of %d", numEvents)
		c.Events = numEvents
	}

	if util.IsStringEmpty(c.LogFile) {
		return fmt.Errorf("empty log file")
	}

	return nil
}

func Override(oldCfg, newCfg *Config) error {
	ov := reflect.ValueOf(oldCfg).Elem()
	nv := reflect.ValueOf(newCfg).Elem()

	overrideFields(ov, nv)

	return nil
}

func overrideFields(ov, nv reflect.Value) {
	for i := 0; i < ov.NumField(); i++ {
		ovField := ov.Field(i)
		if !ovField.CanInterface() {
			continue
		}

		nvField := nv.Field(i)

		if nvField.Kind() == reflect.Struct {
			overrideFields(ovField, nvField)
		} else {
			fv := nvField.Interface()
			isZero := reflect.ValueOf(fv).IsZero()

			if isZero {
				continue
			}

			ovField.Set(reflect.ValueOf(fv))
		}
	}
}
