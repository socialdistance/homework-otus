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
	End         time.Time
	Description string
	UserID      uuid.UUID
}

func NewEvent(title string, started, end time.Time, description string, userID uuid.UUID) *Event {
	return &Event{
		ID:          uuid.New(),
		Title:       title,
		Started:     started,
		End:         end,
		Description: description,
		UserID:      userID,
	}
}
