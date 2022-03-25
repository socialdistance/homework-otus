package app

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/socialdistance/hw12_13_14_15_calendar/internal/logger"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestAppEvent(t *testing.T) {
	userID := uuid.New()
	start, err := time.Parse("2006-01-02 15:04:05", "2022-03-13 12:00:00")
	if err != nil {
		t.FailNow()
		return
	}
	end, err := time.Parse("2006-01-02 15:04:05", "2022-03-14 12:00:00")
	if err != nil {
		t.FailNow()
		return
	}

	logFile, err := os.CreateTemp("", "log")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	logg, err := internallogger.New(config.LoggerConf{
		Level:    config.Info,
		Filename: logFile.Name(),
	})
	if err != nil {
		log.Fatalf("Failed logger %s", err)
	}

	memmoryStorage := memorystorage.New()

	ctx := context.Background()

	testApp := New(logg, memmoryStorage)
	t.Run("Create event", func(t *testing.T) {
		event := storage.Event{
			ID:          uuid.New(),
			Title:       "Test title",
			Started:     start,
			Ended:       end,
			Description: "Test description",
			UserID:      userID,
		}

		err = testApp.CreateEvent(ctx, event)
		require.Nil(t, err)
	})

	t.Run("Update event", func(t *testing.T) {
		event := storage.Event{
			ID:          uuid.New(),
			Title:       "Test title",
			Started:     start,
			Ended:       end,
			Description: "Test description",
			UserID:      userID,
		}

		err = testApp.CreateEvent(ctx, event)
		require.Nil(t, err)

		event.Title = "New title event"
		err = testApp.UpdateEvent(ctx, event)
		require.Nil(t, err)
	})

	t.Run("Delete events", func(t *testing.T) {
		event := storage.Event{
			ID:          uuid.New(),
			Title:       "Test title",
			Started:     start,
			Ended:       end,
			Description: "Test description",
			UserID:      userID,
		}

		err = testApp.CreateEvent(ctx, event)
		require.Nil(t, err)

		err = testApp.DeleteEvent(ctx, event.ID)
		require.Nil(t, err)
	})

	t.Run("Find all events", func(t *testing.T) {
		res, err := testApp.FindAllEvent(ctx)
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, res, 2)
	})
}
