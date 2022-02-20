package system

import (
	"context"
	"net/http"

	"github.com/frain-dev/immune/callback"
	"github.com/frain-dev/immune/net"
	"github.com/pkg/errors"
)

// TODO: more work on the context timeout and cancellations

func (s *System) Run(ctx context.Context) error {
	cfg := callback.Config{
		Port:  80,
		Route: "/",
	}
	var cs *callback.Server
	var err error

	if s.needsCallback {
		cs, err = callback.NewServer(cfg)
		if err != nil {
			return errors.Wrap(err, "failed to initialize new callback se")
		}

		err = cs.Start(ctx)
		if err != nil {
			return err
		}

		defer cs.Stop()
	}

	n := net.NewNet(cs, http.DefaultClient, s.Variables, s.MaxCallbackWaitSeconds, s.BaseURL)

	for i := range s.SetupTestCases {
		err = n.ExecuteSetupTestCase(ctx, &s.SetupTestCases[i])
		if err != nil {
			return err
		}
	}

	for i := range s.TestCases {
		err = n.ExecuteTestCase(ctx, &s.TestCases[i])
		if err != nil {
			return err
		}
	}

	return nil
}
