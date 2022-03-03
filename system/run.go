package system

import (
	"context"
	"net/http"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/callback"
	"github.com/frain-dev/immune/database"
	"github.com/frain-dev/immune/exec"
	"github.com/frain-dev/immune/funcs"
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

	ex := exec.NewExecutor(cs, http.DefaultClient, s.Variables, s.Callback.MaxWaitSeconds, s.BaseURL, s.Callback.IDLocation, truncator)

	//log.Info("starting execution of setup test cases")
	//for i := range s.SetupTestCases {
	//	err = ex.ExecuteSetupTestCase(ctx, &s.SetupTestCases[i])
	//	if err != nil {
	//		return err
	//	}
	//}
	//log.Info("finished execution of setup test cases")

	log.Info("starting execution of test cases")

	for i := range s.TestCases {
		tc := &s.TestCases[i]
		for _, setupName := range tc.Setup {
			switch setupName {
			case "setup_group":
				err = funcs.SetupGroup(ctx, ex)
				if err != nil {
					return err
				}
			case "setup_app":
				err = funcs.SetupApp(ctx, ex)
				if err != nil {
					return err
				}
			case "setup_endpoint":
				err = funcs.SetupAppEndpoint(ctx, s.EventTargetURL, ex)
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
