package callback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/frain-dev/immune"
	log "github.com/sirupsen/logrus"
)

// server is a callback server. It will listen for requests on its
// specified port and parse incoming requests into a *Signal. The
// resulting *Signal is sent on it's outbound channel.
// A callback should be received with it's ReceiveCallback method.
type server struct {
	// all callbacks are sent on this channel
	outbound chan *immune.Signal

	// this is a signal channel, it is never sent on, but once
	// closed by Stop, it will trigger a graceful shutdown of the server
	// see Start
	stop chan struct{}

	// the callback http server
	s *http.Server

	withSSL     bool
	sslKeyFile  string
	sslCertFile string
}

// NewServer instantiates a new callback server
func NewServer(cfg *immune.CallbackConfiguration) (immune.CallbackServer, error) {
	outbound := make(chan *immune.Signal)

	mux := http.DefaultServeMux
	mux.HandleFunc(cfg.Route, handleCallback(outbound))

	srv := &http.Server{
		Addr:    ":" + strconv.FormatUint(uint64(cfg.Port), 10),
		Handler: mux,
	}

	s := &server{
		stop:     make(chan struct{}),
		outbound: outbound,
		s:        srv,
	}

	if cfg.SSL {
		s.withSSL = true
		s.sslKeyFile = cfg.SSLKeyFile
		s.sslCertFile = cfg.SSLCertFile
	}

	return s, nil
}

// Start starts the callback http server
func (s *server) Start(ctx context.Context) error {
	go func() {
		var err error

		if s.withSSL {
			log.Infof("Started callback server with SSL: cert_file: %s, key_file: %s", s.sslCertFile, s.sslKeyFile)
			err = s.s.ListenAndServeTLS(s.sslCertFile, s.sslKeyFile)
		} else {
			err = s.s.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("callback server failed to start")
		}
	}()

	// watches for context cancellation & the stop channel being closed
	go func() {
		select {
		case <-ctx.Done():
			s.gracefulShutdown()
		case <-s.stop:
			s.gracefulShutdown()
		}
	}()
	time.Sleep(time.Second) // allow for the server to start
	return nil
}

// handleCallback returns a http.HandlerFunc that handles a request
// to the callback server
func handleCallback(outbound chan<- *immune.Signal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sig := &immune.Signal{}
		clone := r.Clone(context.Background())

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			sig.Err = fmt.Errorf("failed to read callback body: %v", err)
		} else {
			err = json.Unmarshal(buf, sig)
			if err != nil {
				sig.Err = fmt.Errorf("failed to decode callback body: %v", err)
			}
		}

		clone.Body = io.NopCloser(bytes.NewBuffer(buf))
		sig.Request = clone
		w.WriteHeader(http.StatusOK)
		outbound <- sig
	}
}

// Stop closes the stop channel, which signals a graceful shutdown of the server
// see Start.
func (s *server) Stop() {
	close(s.stop)
}

func (s *server) gracefulShutdown() {
	cctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.s.Shutdown(cctx)
	if err != nil {
		log.WithError(err).Fatal("failed to shutdown callback server")
	}
	log.Infof("callback server shutdown gracefully")
}

// ReceiveCallback sends a Signal to rc
func (s *server) ReceiveCallback(rc chan<- *immune.Signal) {
	rc <- <-s.outbound
}
