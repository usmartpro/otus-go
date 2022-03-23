package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Error(format string, params ...interface{})
}

type Storage interface {
	Select() ([]storage.Event, error)
	Insert(e storage.Event) error
	Update(e storage.Event) error
	Delete(id uuid.UUID) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
