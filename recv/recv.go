package recv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/frain-dev/convoy/pkg/verifier"
	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/auth"
	"github.com/frain-dev/immune/config"
	log "github.com/sirupsen/logrus"
)

type Receiver struct {
	authenticator auth.Authenticator
	cfg           *config.Config
	s             *http.Server
	l             *Log
}

func NewReceiver(cfg *config.Config) *Receiver {
	port := cfg.RecvPort
	if port == 0 {
		port = 80
	}

	rc := &Receiver{
		cfg:           cfg,
		authenticator: auth.NewAuthenticator(&cfg.EndpointConfig),
		s: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		l: NewLog(),
	}

	router := chi.NewRouter()
	router.Post("/", rc.OK)

	rc.s.Handler = router
	return rc
}

func (rc *Receiver) OK(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	rc.l.CaptureHeaders(r, &now)

	if !rc.authenticator.Authenticate(r) {
		rc.l.AddAuthFailure()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := rc.VerifyRequest(r)
	if err != nil {
		rc.l.AddSignatureFailure()
		log.WithError(err).Error("failed to verify request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *Receiver) VerifyRequest(r *http.Request) error {
	opts := &verifier.HmacOptions{
		Header:   immune.DefaultSignatureHeader,
		Hash:     immune.DefaultHash,
		Secret:   rc.cfg.EndpointConfig.Secret,
		Encoding: immune.DefaultEncoding,
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}

	err = verifier.NewHmacVerifier(opts).VerifyRequest(r, payload)
	if err != nil {
		return fmt.Errorf("failed to verify request: %v", err)
	}
	return nil
}

func (rc *Receiver) Listen() *Log {
	go func() {
		// service connections
		if err := rc.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to listen")
		}
	}()

	timeoutChan := make(chan struct{})

	go func() {
		t := time.Duration(rc.cfg.RecvTimeout) * time.Minute
		if t == 0 {
			t = time.Minute * 20
		}
		time.Sleep(t)
		close(timeoutChan)
		fmt.Println("Receive timeout elapsed, shutting down server...")
	}()

	log.Infof("Recv server started")

	rc.gracefulShutdown(timeoutChan)

	rc.l.CalculateStats()
	return rc.l
}

func (rc *Receiver) gracefulShutdown(timeoutChan chan struct{}) {
	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

l:
	for {
		select {
		case <-quit:
			break l
		case <-timeoutChan:
			break l
		}
	}

	log.Info("Stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rc.s.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server Shutdown")
	}

	log.Info("Server exiting")

	time.Sleep(2 * time.Second) // allow all pending connections close themselves
}
