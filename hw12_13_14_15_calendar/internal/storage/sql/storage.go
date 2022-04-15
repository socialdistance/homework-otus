package sqlstorage

import (
	"context"
	"errors"
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

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database %s", err)
	}

	s.conn = conn
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Create(evt storage.Event) error {
	sql := `
		INSERT INTO events (id, title, started, ended, description, user_id) VALUES 
		($1, $2, $3, $4, $5, $6)
	`

	_, err := s.conn.Exec(s.ctx, sql, evt.ID.String(), evt.Title, evt.Started.Format(time.RFC3339),
		evt.Ended.Format(time.RFC3339), evt.Description, evt.UserID)
	return err
}

func (s *Storage) Update(evt storage.Event) error {
	sql := `
		UPDATE events SET title=$2, started=$3, ended=$4, description=$5, user_id=$6 WHERE id=$1
	`

	_, err := s.conn.Exec(s.ctx, sql, evt.ID.String(), evt.Title, evt.Started.Format(time.RFC3339),
		evt.Ended.Format(time.RFC3339), evt.Description, evt.UserID)

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := `
		DELETE FROM events where id=$1
	`

	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) Find(id uuid.UUID) (*storage.Event, error) {
	var event storage.Event
	sql := `select id, title, started_at, finished_at, description, user_id, notify from events where id = $1`

	err := s.conn.QueryRow(s.ctx, sql, id).Scan(
		&event.ID,
		&event.Title,
		&event.Started,
		&event.Ended,
		&event.Description,
		&event.UserID,
		&event.Notify,
	)

	if err == nil {
		return &event, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, fmt.Errorf("cant scan SQL result to struct %w", err)
}

func (s *Storage) FindAll() ([]storage.Event, error) {
	var events []storage.Event

	sql := `
		SELECT id, title, started, ended, description, user_id FROM events
	`

	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var evt storage.Event
		if err := rows.Scan(&evt.ID, &evt.Title, &evt.Started, &evt.Ended, &evt.Description, &evt.UserID); err != nil {
			return nil, fmt.Errorf("cant convert result: %w", err)
		}

		events = append(events, evt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
