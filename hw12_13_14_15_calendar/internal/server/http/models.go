package internalhttp

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Dto struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	StartedAt   string `json:"startedAt"`
	FinishedAt  string `json:"finishedAt"`
	Description string `json:"description"`
	NotifyAt    string `json:"notifyAt"`
}

type ErrorDto struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func ModelToEventDto(event storage.Event) Dto {
	Dto := Dto{}
	Dto.ID = event.ID.String()
	Dto.UserID = event.UserID.String()
	Dto.Title = event.Title
	Dto.StartedAt = event.StartedAt.Format(time.RFC3339)
	Dto.FinishedAt = event.FinishedAt.Format(time.RFC3339)
	Dto.Description = event.Description
	Dto.NotifyAt = event.NotifyAt.Format(time.RFC3339)

	return Dto
}

func (d *Dto) GetModel() (*storage.Event, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, fmt.Errorf("ID must be UUID, your value: %s, %w", d.ID, err)
	}

	userID, err := uuid.Parse(d.UserID)
	if err != nil {
		return nil, fmt.Errorf("userID must be UUID, your value: %s, %w", d.UserID, err)
	}

	startedAt, err := time.Parse("2006-01-02 15:04:00", d.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("startedAt format must be 'yyyy-mm-dd hh:mm:ss', your value: %s, %w", d.StartedAt, err)
	}

	finishedAt, err := time.Parse("2006-01-02 15:04:00", d.FinishedAt)
	if err != nil {
		return nil, fmt.Errorf("finishedAt format must be 'yyyy-mm-dd hh:mm:ss', your value: %s, %w", d.FinishedAt, err)
	}

	notifyAt, err := time.Parse("2006-01-02 15:04:00", d.NotifyAt)
	if err != nil {
		return nil, fmt.Errorf("notifyAt format must be 'yyyy-mm-dd hh:mm:ss', your value: %s, %w", d.NotifyAt, err)
	}

	appEvent := storage.NewEvent(userID, d.Title, startedAt, finishedAt, d.Description, notifyAt)
	appEvent.ID = id

	return appEvent, nil
}
