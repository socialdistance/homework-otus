package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDateBusy     = errors.New("error date busy")
	ErrDateNotExist = errors.New("error date not exists")
)

type Event struct {
	ID          uuid.UUID
	Title       string
	Started     time.Time
	Ended       time.Time
	Description string
	UserID      uuid.UUID
}

func NewEvent(title string, started, ended time.Time, description string, userID uuid.UUID) *Event {
	return &Event{
		ID:          uuid.New(),
		Title:       title,
		Started:     started,
		Ended:       ended,
		Description: description,
		UserID:      userID,
	}
}
