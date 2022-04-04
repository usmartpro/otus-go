package memorystorage

import (
	"sync"

	"github.com/google/uuid"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) Insert(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		return storage.ErrDuplicateEvent
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Update(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrNotExistEvent
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) Select() ([]storage.Event, error) {
	events := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) SelectOne(id uuid.UUID) (*storage.Event, error) {
	if event, ok := s.events[id]; ok {
		return &event, nil
	}

	return nil, nil
}
