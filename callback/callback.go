package callback

import (
	"context"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	outbound chan struct{}
	stop     chan struct{}
	s        *http.Server
}

type Config struct {
	Port  uint
	Route string
}

func NewServer(cfg Config) (*Server, error) {
	outbound := make(chan struct{})

	mux := http.DefaultServeMux
	mux.HandleFunc(cfg.Route, func(w http.ResponseWriter, r *http.Request) {
		outbound <- struct{}{}
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
}

func (s *Server) ReceiveCallback() struct{} {
	return <-s.outbound
}
