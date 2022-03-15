package sqlstorage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	ctx  context.Context
	conn *pgx.Conn
	url  string
}

func New(ctx context.Context, url string) *Storage {
	return &Storage{
		ctx: ctx,
		url: url,
	}
}

func (s *Storage) Connect(ctx context.Context, url string) error {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect database %s", err)
		os.Exit(1)
	}

	s.conn = conn
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Create(evt storage.Event) error {
	sql := `
		INSERT INTO events (id, title, started, end, description, user_id) VALUES 
		($1, $2, $3, $4, $5, $6)
	`

	_, err := s.conn.Exec(s.ctx, sql, evt.ID.String(), evt.Title, evt.Started.Format(time.RFC3339),
		evt.End.Format(time.RFC3339), evt.Description, evt.UserID)
	return err
}

func (s *Storage) Update(evt storage.Event) error {
	sql := `
		UPDATE events SET title=$2, started=$3, end=$4, description=$5 where id=$1
	`

	_, err := s.conn.Exec(s.ctx, sql, evt.ID.String(), evt.Title, evt.Started.Format(time.RFC3339),
		evt.End.Format(time.RFC3339), evt.Description, evt.UserID)

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := `
		DELETE FROM events where id=$1
	`

	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) FindAll() ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	sql := `
		SELECT id, title, started, end, description, user_id from events ORDER BY date
	`

	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var evt storage.Event
		if err := rows.Scan(&evt.ID, &evt.Title, &evt.Started, &evt.End, &evt.Description, &evt.UserID); err != nil {
			return nil, fmt.Errorf("cant convert result: %w", err)
		}

		events = append(events, evt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
