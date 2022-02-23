package callback

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	outbound chan *Signal
	stop     chan struct{}
	s        *http.Server
}

type Config struct {
	Port  uint
	Route string
}

type Signal struct {
	ImmuneCallBackID string `json:"immune_callback_id"`
}

func NewServer(cfg Config) (*Server, error) {
	outbound := make(chan *Signal)

	mux := http.DefaultServeMux
	mux.HandleFunc(cfg.Route, func(w http.ResponseWriter, r *http.Request) {
		sig := &Signal{}
		err := json.NewDecoder(r.Body).Decode(sig)
		if err != nil {
			log.WithError(err).Error("failed to decode callback body")
			return
		}
		outbound <- sig
	})

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

func (s *Server) Start(ctx context.Context) error {
	go func() {
		err := s.s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("callback server failed to start")
		}
	}()

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

func (s *Server) ReceiveCallback() *Signal {
	return <-s.outbound
}
