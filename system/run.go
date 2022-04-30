package system

import (
	"context"
	"fmt"
	"net/http"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/callback"
	"github.com/frain-dev/immune/database"
	"github.com/frain-dev/immune/exec"
	"github.com/frain-dev/immune/funcs"
	"github.com/google/uuid"
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

	truncator, err := database.NewTruncator(&s.Database)
	if err != nil {
		return err
	}

	err = truncator.Truncate(ctx)
	if err != nil {
		return err
	}

	idFn := func() string {
		return uuid.New().String()
	}

	sv, err := callback.NewSignatureVerifier(
		s.Callback.Signature.ReplayAttacks,
		s.Callback.Signature.Secret,
		s.Callback.Signature.Header,
		s.Callback.Signature.Hash,
	)
	if err != nil {
		return fmt.Errorf("failed to get new signature verifier: %v", err)
	}

	ex := exec.NewExecutor(
		cs,
		http.DefaultClient,
		s.Variables,
		sv,
		s.Callback.MaxWaitSeconds,
		s.BaseURL,
		s.Callback.IDLocation,
		truncator,
		idFn,
	)

	log.Info("starting execution of test cases")
	for i := range s.TestCases {
		tc := &s.TestCases[i]
		for _, setupName := range tc.Setup {
			switch setupName {
			case "setup_group":
				err = funcs.SetupGroup(ctx, ex, &s.Callback.Signature)
				if err != nil {
					return err
				}
			case "setup_app":
				err = funcs.SetupApp(ctx, ex)
				if err != nil {
					return err
				}
			case "setup_endpoint":
				err = funcs.SetupAppEndpoint(ctx, s.EventTargetURL, s.Callback.Signature.Secret, ex)
				if err != nil {
					return err
				}
			case "setup_event":
				err = funcs.SetupEvent(ctx, ex)
				if err != nil {
					return err
				}
			default:
				return errors.Errorf("unknown setup %s, in test case %s", setupName, tc.Name)
			}
		}
		err = ex.ExecuteTestCase(ctx, tc)
		if err != nil {
			return err
		}
		log.Infof("test_case %s passed", tc.Name)
	}

	log.Info("finished execution of test cases")

	return nil
}
