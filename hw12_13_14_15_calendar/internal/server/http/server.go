package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	host   string
	port   string
	logger Logger
	server *http.Server
}

type Logger interface {
	Debug(message string, params ...interface{})
	Info(message string, params ...interface{})
	Error(message string, params ...interface{})
	Warn(message string, params ...interface{})
	LogHTTP(r *http.Request, code, length int)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app *app.App, host, port string) *Server {
	server := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	httpServ := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: loggingMiddleware(http.HandlerFunc(server.HandleHTTP), logger),
	}

	server.server = httpServ

	return server
}

func Routers(app *app.App) http.Handler {
	handlers := NewServerHandlers(app)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HelloWorld).Methods("GET")
	r.HandleFunc("/create", handlers.CreateEvent).Methods("POST")
	r.HandleFunc("/events/update/{id}", handlers.UpdateEvent).Methods("PUT")
	r.HandleFunc("/events/delete/{id}", handlers.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events", handlers.ListEvents).Methods("GET")

	return r
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("[+] Http server start and listen %s:%s", s.host, s.port)
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
