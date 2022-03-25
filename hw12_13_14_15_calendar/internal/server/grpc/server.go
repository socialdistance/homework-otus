//go:generate protoc --go_out=. --go-grpc_out=. ../../../api/EventService.proto --proto_path=../../../api

package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	internalapp "github.com/socialdistance/hw12_13_14_15_calendar/internal/app"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedEventServiceServer
	host       string
	port       string
	logger     Logger
	grpcServer *grpc.Server
	app        *internalapp.App
}

type Logger interface {
	Debug(message string, params ...interface{})
	Info(message string, params ...interface{})
	Error(message string, params ...interface{})
	Warn(message string, params ...interface{})
	LogHTTP(r *http.Request, code, length int)
}

type Application interface{}

func NewServerLogger(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		logger.Info("GRPC Request: %v", req)
		return handler(ctx, req)
	}
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(NewServerLogger(logger)),
	)

	server := &Server{
		host:       host,
		port:       port,
		logger:     logger,
		grpcServer: grpcServer,
	}

	RegisterEventServiceServer(server.grpcServer, server)

	return server
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("[-] Can't start GRPC server: %w", err)
	}

	s.logger.Info("[+] Start GRPC Server %s:%s", s.host, s.port)

	err = s.grpcServer.Serve(lsn)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) Create(ctx context.Context, in *Event) (*ResponseEvent, error) {
	appEvt := storage.Event{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid ID value. Except UUID, got %s, %w", in.GetId(), err)
	}
	appEvt.ID = id

	userID, err := uuid.Parse(in.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("invalid userID value. Except userID, got %s, %w", in.GetUserID(), err)
	}
	appEvt.UserID = userID

	started, err := time.Parse("2006-01-02 15:04:05", in.GetStarted())
	if err != nil {
		return nil, fmt.Errorf("invalid started value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvt.Started = started

	ended, err := time.Parse("2006-01-02 15:04:05", in.GetStarted())
	if err != nil {
		return nil, fmt.Errorf("invalid ended value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvt.Ended = ended

	fmt.Println(appEvt)

	err = s.app.CreateEvent(ctx, appEvt)
	if err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) Delete(ctx context.Context, in *DeleteEvent) (*ResponseEvent, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid ID value. Except UUID, got %s, %w", in.GetId(), err)
	}

	err = s.app.DeleteEvent(ctx, id)
	if err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) Update(ctx context.Context, in *Event) (*ResponseEvent, error) {
	appEvt := storage.Event{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid ID value. Except UUID, got %s, %w", in.GetId(), err)
	}
	appEvt.ID = id

	userID, err := uuid.Parse(in.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("invalid userID value. Except userID, got %s, %w", in.GetUserID(), err)
	}
	appEvt.UserID = userID

	started, err := time.Parse("2006-01-02 15:04:05", in.GetStarted())
	if err != nil {
		return nil, fmt.Errorf("invalid started value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvt.Started = started

	ended, err := time.Parse("2006-01-02 15:04:05", in.GetStarted())
	if err != nil {
		return nil, fmt.Errorf("invalid ended value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvt.Ended = ended

	if err = s.app.UpdateEvent(ctx, appEvt); err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func ListResponse(evts []storage.Event) *ResponseEventList {
	response := ResponseEventList{}
	for _, t := range evts {
		response.Event = append(response.Event, &Event{
			Id:          t.ID.String(),
			Title:       t.Title,
			Started:     t.Started.Format(time.RFC3339),
			Ended:       t.Ended.Format(time.RFC3339),
			Description: t.Description,
			UserID:      t.UserID.String(),
		})
	}

	return &response
}

func ResponseSuccess() *ResponseEvent {
	return &ResponseEvent{
		Success: true,
		Error:   "",
	}
}

func ResponseError(msg string) *ResponseEvent {
	return &ResponseEvent{
		Success: false,
		Error:   msg,
	}
}
