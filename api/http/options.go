package http

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func ServerAddr(addr string) func(server *Server) {
	return func(server *Server) {
		server.server.Addr = addr
	}
}

func ServerRouter(r chi.Router) func(server *Server) {
	return func(server *Server) {
		server.router = r
		server.server.Handler = server.router
	}
}

func ServerLogger(l logrus.FieldLogger) func(server *Server) {
	return func(server *Server) {
		server.logger = l
	}
}

func ServerReadTimout(duration time.Duration) func(server *Server) {
	return func(server *Server) {
		server.server.ReadTimeout = duration
	}
}
func ServerWriteTimeout(duration time.Duration) func(server *Server) {
	return func(server *Server) {
		server.server.WriteTimeout = duration
	}
}
