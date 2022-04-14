//go:generate protoc --go_out=. --go-grpc_out=. ../../../api/EventService.proto --proto_path=../../../api

package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	internalapp "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedEventServiceServer
	host    string
	port    string
	grpcSrv *grpc.Server
	app     *internalapp.App
	logg    internalapp.Logger
}

func NewServerLogger(logger internalapp.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		logger.Info("NEW GRPC Request: %v", req)
		return handler(ctx, req)
	}
}

func NewServer(logger internalapp.Logger, app *internalapp.App, host string, port string) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			NewServerLogger(logger),
		),
	)

	s := &Server{
		host:    host,
		port:    port,
		grpcSrv: grpcServer,
		app:     app,
		logg:    logger,
	}

	RegisterEventServiceServer(s.grpcSrv, s)

	return s
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	s.logg.Info("HTTP server run %s:%s", s.host, s.port)

	return s.grpcSrv.Serve(lsn)
}

func (s *Server) Stop() {
	s.grpcSrv.GracefulStop()
}

func (s *Server) Create(ctx context.Context, in *Event) (*EventResponse, error) { // nolint:dupl
	appEvent := storage.Event{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid ID value. Exprected UUID, got %s, %w", in.GetId(), err)
	}
	appEvent.ID = id

	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("invalid UserID value. Exprected UUID, got %s, %w", in.GetId(), err)
	}
	appEvent.UserID = userID

	startedAt, err := time.Parse("2006-01-02 15:04:05", in.GetStartedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid StartedAt value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvent.StartedAt = startedAt

	finishedAt, err := time.Parse("2006-01-02 15:04:05", in.GetStartedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid FinishedAt value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvent.FinishedAt = finishedAt

	notifyAt, err := time.Parse("2006-01-02 15:04:05", in.GetStartedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid Notify value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvent.NotifyAt = notifyAt

	if err = s.app.CreateEvent(ctx, appEvent); err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) Update(ctx context.Context, in *Event) (*EventResponse, error) { // nolint:dupl
	appEvent := storage.Event{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
	}

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid ID value. Exprected UUID, got %s, %w", in.GetId(), err)
	}
	appEvent.ID = id

	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("invalid UserID value. Exprected UUID, got %s, %w", in.GetId(), err)
	}
	appEvent.UserID = userID

	startedAt, err := time.Parse("2006-01-02 15:04:05", in.GetStartedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid StartedAt value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvent.StartedAt = startedAt

	finishedAt, err := time.Parse("2006-01-02 15:04:05", in.GetStartedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid FinishedAt value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvent.FinishedAt = finishedAt

	notifyAt, err := time.Parse("2006-01-02 15:04:05", in.GetStartedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid Notify value. Exprected 2006-01-02 15:04:05, got %s, %w", in.GetId(), err)
	}
	appEvent.NotifyAt = notifyAt

	if err = s.app.UpdateEvent(ctx, appEvent); err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) Delete(ctx context.Context, in *DeleteEventRequest) (*EventResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid ID value. Exprected UUID, got %s,%w", in.GetId(), err)
	}

	if err = s.app.DeleteEvent(ctx, id); err != nil {
		return ResponseError(err.Error()), nil
	}

	return ResponseSuccess(), nil
}

func (s *Server) DayEvents(ctx context.Context, in *EventsRequest) (*EventsResponse, error) {
	dtStart, err := time.Parse("2006-01-02", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Expected yyyy-mm-dd, got %s", in.GetDate())
	}

	events, err := s.app.GetDayEvents(ctx, dtStart)
	if err != nil {
		return nil, err
	}

	return ListResponseSuccess(events), nil
}

func (s *Server) WeekEvents(ctx context.Context, in *EventsRequest) (*EventsResponse, error) {
	dtStart, err := time.Parse("2006-01-02", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Expected yyyy-mm-dd, got %s", in.GetDate())
	}

	events, err := s.app.GetWeekEvents(ctx, dtStart)
	if err != nil {
		return nil, err
	}

	return ListResponseSuccess(events), nil
}

func (s *Server) MonthEvents(ctx context.Context, in *EventsRequest) (*EventsResponse, error) {
	dtStart, err := time.Parse("2006-01-02", in.GetDate())
	if err != nil {
		return nil, fmt.Errorf("invalid date value. Expected yyyy-mm-dd, got %s", in.GetDate())
	}

	events, err := s.app.GetMonthEvents(ctx, dtStart)
	if err != nil {
		return nil, err
	}

	return ListResponseSuccess(events), nil
}

func ResponseSuccess() *EventResponse {
	return &EventResponse{
		Result: true,
		Error:  "",
	}
}

func ListResponseSuccess(events []storage.Event) *EventsResponse {
	resp := EventsResponse{}
	for _, e := range events {
		resp.Events = append(resp.Events, &Event{
			Id:          e.ID.String(),
			UserId:      e.UserID.String(),
			Title:       e.Title,
			StartedAt:   e.StartedAt.Format(time.RFC3339),
			FinishedAt:  e.StartedAt.Format(time.RFC3339),
			Description: e.Description,
			NotifyAt:    e.NotifyAt.Format(time.RFC3339),
		})
	}

	return &resp
}

func ResponseError(msg string) *EventResponse {
	return &EventResponse{
		Result: false,
		Error:  msg,
	}
}
