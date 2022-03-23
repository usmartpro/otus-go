package internalhttp

import (
	"context"
	"net"
	"net/http"
)

type Server struct {
	host   string
	port   string
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(message string, params ...interface{})
	Error(message string, params ...interface{})
	LogRequest(r *http.Request, code, length int)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	httpServer := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	newServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: loggingMiddleware(http.HandlerFunc(httpServer.handleHTTP), logger),
	}

	httpServer.server = newServer

	return httpServer
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello-world"))
	w.WriteHeader(200)
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("HTTP server run %s:%s", s.host, s.port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
