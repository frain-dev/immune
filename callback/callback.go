package callback

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/frain-dev/immune"

	log "github.com/sirupsen/logrus"
)

// Server is a callback server. It will listen for requests on its
// specified port and parse incoming requests into a *Signal. The
// resulting *Signal is sent on it's outbound channel.
// A callback should be received with it's ReceiveCallback method.
type Server struct {
	// all callback are sent on this channel
	outbound chan *Signal

	// this is a signal channel, it is never sent on, but once
	// closed by Stop, it will trigger a graceful shutdown of the server
	stop chan struct{}

	// the callback http.Server
	s *http.Server
}

// A Signal represents a single callback
type Signal struct {
	// ImmuneCallBackID collects the callback id from the request body, it's json tag
	// must always match immune.CallbackIDFieldName
	ImmuneCallBackID string `json:"immune_callback_id"`
}

// NewServer instantiates a new callback Server
func NewServer(cfg *immune.CallbackConfiguration) (*Server, error) {
	outbound := make(chan *Signal)

	mux := http.DefaultServeMux
	mux.HandleFunc(cfg.Route, handleCallback(outbound))

	srv := &http.Server{
		Addr:    ":" + strconv.FormatUint(uint64(cfg.Port), 10),
		Handler: mux,
	}

	return &Server{
		stop:     make(chan struct{}),
		outbound: outbound,
		s:        srv,
	}, nil
}

// Start starts the callback http server
func (s *Server) Start(ctx context.Context) error {
	go func() {
		err := s.s.ListenAndServe()
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

	return nil
}

// handleCallback returns a http.HandlerFunc that handles a request
// to the callback server
func handleCallback(outbound chan *Signal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sig := &Signal{}
		err := json.NewDecoder(r.Body).Decode(sig)
		if err != nil {
			log.WithError(err).Error("failed to decode callback body")
			return
		}
		outbound <- sig
	}
}

// Stop closes the stop channel, which signals a graceful shutdown of the server
// see Start.
func (s *Server) Stop() {
	close(s.stop)
}

func (s *Server) gracefulShutdown() {
	cctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.s.Shutdown(cctx)
	if err != nil {
		log.WithError(err).Fatal("failed to shutdown callback server")
	}
	log.Infof("callback server shutdown gracefully")
}

// ReceiveCallback receives a Signal from the callback channel
func (s *Server) ReceiveCallback() *Signal {
	return <-s.outbound
}
