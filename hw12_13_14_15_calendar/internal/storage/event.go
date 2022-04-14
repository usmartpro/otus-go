package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDuplicateEvent = errors.New("error duplicate event")
	ErrNotExistEvent  = errors.New("event not exist")
)

type Event struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Title            string
	StartedAt        time.Time
	FinishedAt       time.Time
	Description      string
	NotifyBeforeTime time.Time
}

func NewEvent(
	userID uuid.UUID,
	title string,
	startedAt time.Time,
	finishedAt time.Time,
	description string,
	notifyBeforeTime time.Time,
) *Event {
	return &Event{
		ID:               uuid.New(),
		UserID:           userID,
		Title:            title,
		StartedAt:        startedAt,
		FinishedAt:       finishedAt,
		Description:      description,
		NotifyBeforeTime: notifyBeforeTime,
	}
}
