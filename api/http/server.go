package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
	router chi.Router
	logger logrus.FieldLogger
}

func NewServer(options ...func(s *Server)) *Server {
	srv := &Server{server: &http.Server{}}
	for _, option := range options {
		option(srv)
	}
	return srv
}

var ErrServerClosed = http.ErrServerClosed

func (s *Server) Post(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.Post(pattern, f)
}

func (s *Server) Put(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.Put(pattern, f)
}

func (s *Server) Get(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.Get(pattern, f)
}

func (s *Server) Delete(pattern string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.Delete(pattern, f)
}

func (s *Server) Start() error {
	s.logger.Infof("# starting http server #")
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
