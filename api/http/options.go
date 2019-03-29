package http

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// ServerAddr adds server listening port
func ServerAddr(addr string) func(server *Server) {
	return func(server *Server) {
		server.server.Addr = addr
	}
}

// ServerRouter adds router to server
func ServerRouter(r chi.Router) func(server *Server) {
	return func(server *Server) {
		server.router = r
		server.server.Handler = server.router
	}
}

// ServerLogger adds logger to server
func ServerLogger(l logrus.FieldLogger) func(server *Server) {
	return func(server *Server) {
		server.logger = l
	}
}

// ServerLogger adds read timeout for server
func ServerReadTimout(duration time.Duration) func(server *Server) {
	return func(server *Server) {
		server.server.ReadTimeout = duration
	}
}

// ServerLogger adds write timeout for server
func ServerWriteTimeout(duration time.Duration) func(server *Server) {
	return func(server *Server) {
		server.server.WriteTimeout = duration
	}
}
