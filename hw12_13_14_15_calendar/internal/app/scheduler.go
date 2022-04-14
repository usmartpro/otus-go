package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type EventSource interface {
	GetActualNotifyEvents(notifyTime time.Time) ([]storage.Event, error)
	GetOldEvents(timeBefore time.Time) ([]storage.Event, error)
	Delete(id uuid.UUID) error
}

type NotificationReceiver interface {
	Add(Notification) error
}

type Scheduler struct {
	eventSource          EventSource
	notificationReceiver NotificationReceiver
	logger               Logger
}

func NewAppScheduler(es EventSource, rcv NotificationReceiver, logger Logger) *Scheduler {
	return &Scheduler{
		es,
		rcv,
		logger,
	}
}

func (s *Scheduler) Notify() error {
	notifyTime := time.Now()
	events, err := s.eventSource.GetActualNotifyEvents(notifyTime)
	if err != nil {
		return fmt.Errorf("error get events for date %s: %w", notifyTime, err)
	}

	if len(events) > 0 {
		s.logger.Info("! scheduler. %d actual messages", len(events))
	} else {
		s.logger.Info("! scheduler. No messages")
	}

	for _, event := range events {
		notification := Notification{
			EventID:  event.ID,
			UserID:   event.UserID,
			Title:    event.Title,
			DateTime: event.NotifyAt,
		}

		if err := s.notificationReceiver.Add(notification); err != nil {
			return fmt.Errorf("error add notification for event %s:  %w", event.ID, err)
		}

		s.logger.Info("! scheduler. Event id=%s sent", notification.EventID)
	}

	return nil
}

func (s *Scheduler) RemoveOldEvents() error {
	timeBefore := time.Now().AddDate(-1, 0, 0)

	events, err := s.eventSource.GetOldEvents(timeBefore)
	if err != nil {
		return fmt.Errorf("error get events for date %s: %w", timeBefore, err)
	}

	for _, event := range events {
		s.eventSource.Delete(event.ID)

		s.logger.Info("! scheduler. Old Event %s removed. Date: %s", event.ID, event.NotifyAt)
	}

	return nil
}
