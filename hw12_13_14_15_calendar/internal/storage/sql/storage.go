package sqlstorage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	ctx  context.Context
	conn *pgx.Conn
	dsn  string
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{
		ctx: ctx,
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) app.Storage {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connect to database: %v\n", err)
		os.Exit(1)
	}
	s.conn = conn
	return s
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Insert(event storage.Event) error {
	fmt.Print("Storage Insert")
	sql := `INSERT INTO events (id, user_id, title, started_at, finished_at, description, notify_before_time) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		event.ID.String(),
		event.UserID,
		event.Title,
		event.StartedAt.Format(time.RFC3339),
		event.FinishedAt.Format(time.RFC3339),
		event.Description,
		event.NotifyBeforeTime.Format(time.RFC3339),
	)

	return err
}

func (s *Storage) Update(event storage.Event) error {
	sql := `UPDATE events 
			SET
				user_id = $1,
			    title = $2,
    			started_at = $3,
    			finished_at = $4,
    			description = $5,
    			notify_before_time = $6
			WHERE id = $7`

	_, err := s.conn.Exec(
		s.ctx,
		sql,
		event.UserID,
		event.Title,
		event.StartedAt.Format(time.RFC3339),
		event.FinishedAt.Format(time.RFC3339),
		event.Description,
		event.NotifyBeforeTime.Format(time.RFC3339),
		event.ID.String(),
	)

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := "DELETE FROM events WHERE id = $1"
	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) Select() ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	sql := `SELECT id, user_id, title, started_at, finished_at, description, notify_before_time 
			FROM events
			ORDER BY id`

	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.StartedAt,
			&event.FinishedAt,
			&event.Description,
			&event.NotifyBeforeTime,
		); err != nil {
			return nil, fmt.Errorf("error scan result: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
