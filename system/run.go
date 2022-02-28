package system

import (
	"context"
	"net/http"

	"github.com/frain-dev/immune"

	"github.com/frain-dev/immune/callback"
	"github.com/frain-dev/immune/exec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Run executes the entire system, starting with the setup test cases,
// and then the test cases, a callback server will be started if needed.
func (s *System) Run(ctx context.Context) error {
	var cs immune.CallbackServer
	var err error

	if s.needsCallback {
		cs, err = callback.NewServer(&s.Callback)
		if err != nil {
			return errors.Wrap(err, "failed to initialize new callback server")
		}

		err = cs.Start(ctx)
		if err != nil {
			return err
		}

		defer cs.Stop()
	}

	ex := exec.NewExecutor(cs, http.DefaultClient, s.Variables, s.Callback.MaxWaitSeconds, s.BaseURL, s.Callback.IDLocation)

	log.Info("starting execution of setup test cases")

	for i := range s.SetupTestCases {
		err = ex.ExecuteSetupTestCase(ctx, &s.SetupTestCases[i])
		if err != nil {
			return err
		}
	}

	log.Info("finished execution of setup test cases")
	log.Info("starting execution of test cases")

	for i := range s.TestCases {
		err = ex.ExecuteTestCase(ctx, &s.TestCases[i])
		if err != nil {
			return err
		}
	}

	log.Info("finished execution of test cases")

	return nil
}
