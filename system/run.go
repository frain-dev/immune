package system

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/frain-dev/immune/callback"
	"github.com/frain-dev/immune/net"
	"github.com/pkg/errors"
)

// TODO: more work on the context timeout and cancellations

func (s *System) Run(ctx context.Context) error {
	cfg := callback.Config{
		Port:  s.Callback.Port,
		Route: s.Callback.Route,
	}
	var cs *callback.Server
	var err error

	if s.needsCallback {
		cs, err = callback.NewServer(cfg)
		if err != nil {
			return errors.Wrap(err, "failed to initialize new callback server")
		}

		err = cs.Start(ctx)
		if err != nil {
			return err
		}

		defer cs.Stop()
	}

	n := net.NewNet(cs, http.DefaultClient, s.Variables, s.Callback.MaxWaitSeconds, s.BaseURL, s.Callback.IDLocation)

	log.Info("starting execution of setup test cases")

	for i := range s.SetupTestCases {
		err = n.ExecuteSetupTestCase(ctx, &s.SetupTestCases[i])
		if err != nil {
			return err
		}
	}

	log.Info("finished execution of setup test cases")
	log.Info("starting execution of test cases")

	for i := range s.TestCases {
		err = n.ExecuteTestCase(ctx, &s.TestCases[i])
		if err != nil {
			return err
		}
	}

	log.Info("finished execution of test cases")

	return nil
}
