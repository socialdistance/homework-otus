package internalhttp

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/storage"
)

type EventDto struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Started     string `json:"started"`
	Ended       string `json:"ended"`
	Description string `json:"description"`
	UserID      string `json:"userID"` // nolint:tagliatelle
}

type ErrorDto struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (e *EventDto) GetModel() (*storage.Event, error) {
	started, err := time.Parse("2006-01-02 15:04:00", e.Started)
	if err != nil {
		return nil, fmt.Errorf("error: Start exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", e.Started, err)
	}

	ended, err := time.Parse("2006-01-02 15:04:00", e.Ended)
	if err != nil {
		return nil, fmt.Errorf("error: End exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", e.Ended, err)
	}

	id, err := uuid.Parse(e.ID)
	if err != nil {
		return nil, fmt.Errorf("ID exprected to be uuid, got: %s, %w", e.ID, err)
	}

	userID, err := uuid.Parse(e.UserID)
	if err != nil {
		return nil, fmt.Errorf("userID exprected to be uuid, got: %s, %w", e.UserID, err)
	}

	appEvent := storage.NewEvent(e.Title, started, ended, e.Description, userID)
	appEvent.ID = id

	return appEvent, nil
}

func CreateDto(event storage.Event) EventDto {
	eventDto := EventDto{
		ID:          event.ID.String(),
		Title:       event.Title,
		Started:     event.Started.Format(time.RFC3339),
		Ended:       event.Ended.Format(time.RFC3339),
		Description: event.Description,
		UserID:      event.UserID.String(),
	}

	return eventDto
}
