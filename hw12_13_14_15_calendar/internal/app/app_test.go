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
	start, err := time.Parse("2006-01-02 15:04:05", "2022-04-10 12:00:00")
	require.Nil(t, err)

	end, err := time.Parse("2006-01-02 15:04:05", "2022-04-11 12:00:00")
	require.Nil(t, err)

	logFile, err := os.CreateTemp("", "log")
	require.Nil(t, err)

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
		require.Nil(t, err)
		require.Len(t, res, 2)
	})

	t.Run("Find events by interval", func(t *testing.T) {
		eventDay := storage.Event{
			ID:          parseUUID(t, "be00b01a-ef9b-4800-bd18-c07798598d6a"),
			Title:       "Event Title 1",
			Started:     parseTime(t, "2022-04-11T12:30:00Z"),
			Ended:       parseTime(t, "2022-04-12T12:30:00Z"),
			Description: "Event Description 1",
			UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
			Notify:      parseTime(t, "2022-04-10T12:30:00Z"),
		}
		err = testApp.CreateEvent(ctx, eventDay)
		require.Nil(t, err)

		eventWeek := storage.Event{
			ID:          parseUUID(t, "be00b01a-ef9b-4800-bd18-c07798598d8a"),
			Title:       "Event Title 2",
			Started:     parseTime(t, "2022-04-15T12:30:00Z"),
			Ended:       parseTime(t, "2022-04-16T12:30:00Z"),
			Description: "Event Description 2",
			UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
			Notify:      parseTime(t, "2022-04-14T12:30:00Z"),
		}
		err = testApp.CreateEvent(ctx, eventWeek)
		require.Nil(t, err)

		eventMonth := storage.Event{
			ID:          parseUUID(t, "be00b01a-ef9b-4800-bd18-c07798598d5a"),
			Title:       "Event Title 3",
			Started:     parseTime(t, "2022-05-10T12:30:00Z"),
			Ended:       parseTime(t, "2022-05-11T12:30:00Z"),
			Description: "Event Description 3",
			UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
			Notify:      parseTime(t, "2022-05-09T12:30:00Z"),
		}
		err = testApp.CreateEvent(ctx, eventMonth)
		require.Nil(t, err)

		event := storage.Event{
			ID:          parseUUID(t, "be00b01a-ef9b-4800-bd18-c07798598d3a"),
			Title:       "Event Title 4",
			Started:     parseTime(t, "2022-04-10T12:30:00Z"),
			Ended:       parseTime(t, "2022-04-11T12:30:00Z"),
			Description: "Event Description 4",
			UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
			Notify:      parseTime(t, "2022-04-09T12:30:00Z"),
		}
		err = testApp.CreateEvent(ctx, event)
		require.Nil(t, err)

		start, _ := time.Parse(time.RFC3339, "2022-04-11T12:30:00Z")

		events, err := testApp.EventsByDay(ctx, start)
		require.Nil(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "be00b01a-ef9b-4800-bd18-c07798598d6a", events[0].ID.String())

		events, err = testApp.EventsByWeek(ctx, start)
		require.Nil(t, err)
		require.Len(t, events, 2)
		require.Equal(t, "be00b01a-ef9b-4800-bd18-c07798598d6a", events[0].ID.String())
		require.Equal(t, "be00b01a-ef9b-4800-bd18-c07798598d8a", events[1].ID.String())

		events, err = testApp.EventsByMonth(ctx, start)
		require.Nil(t, err)
		require.Len(t, events, 3)
		require.Equal(t, "be00b01a-ef9b-4800-bd18-c07798598d6a", events[0].ID.String())
		require.Equal(t, "be00b01a-ef9b-4800-bd18-c07798598d8a", events[1].ID.String())
		require.Equal(t, "be00b01a-ef9b-4800-bd18-c07798598d5a", events[2].ID.String())
	})
}

func parseUUID(t *testing.T, stringUUID string) uuid.UUID {
	t.Helper()

	id, err := uuid.Parse(stringUUID)
	if err != nil {
		t.Errorf("can't parse uuid %s", err)
	}

	return id
}

func parseTime(t *testing.T, stringTime string) time.Time {
	t.Helper()

	timeDuration, err := time.Parse(time.RFC3339, stringTime)
	if err != nil {
		t.Errorf("can't parse time %s", err)
	}

	return timeDuration
}
