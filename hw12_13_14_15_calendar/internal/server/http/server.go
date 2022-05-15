package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
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

func NewServer(logger Logger, app *app.App, host, port string) *Server {
	httpServer := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	newServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: loggingMiddleware(NewRouter(app), logger),
	}

	httpServer.server = newServer

	return httpServer
}

func NewRouter(app *app.App) http.Handler {
	handlers := NewServerHandlers(app)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HelloWorld).Methods("GET")
	r.HandleFunc("/events", handlers.CreateEvent).Methods("POST")
	r.HandleFunc("/events/{id}", handlers.UpdateEvent).Methods("PUT")
	r.HandleFunc("/events/{id}", handlers.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events", handlers.ListEvents).Methods("GET")

	return r
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
	var err error
	if err = s.server.Shutdown(ctx); err == nil {
		s.logger.Info("HTTP server stopped")
	}
	return err
}
