package memorystorage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	memorystorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := New()

	t.Run("storage test", func(t *testing.T) {
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

		evt := memorystorage.NewEvent(
			"title",
			start,
			end,
			"description",
			userID,
		)

		err = storage.Create(*evt)
		if err != nil {
			t.FailNow()
			return
		}

		res, _ := storage.FindAll()
		require.Len(t, res, 1)
		require.Equal(t, *evt, res[0])

		evt.Title = "New event title"

		res, _ = storage.FindAll()
		require.Len(t, res, 1)
		require.NotEqual(t, *evt, res[0])
		require.NotEqual(t, evt.Title, res[0].Title)

		err = storage.Update(*evt)
		if err != nil {
			t.FailNow()
			return
		}

		res, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, res, 1)
		require.Equal(t, *evt, res[0])
		require.Equal(t, evt.Title, res[0].Title)

		err = storage.Delete(evt.ID)
		if err != nil {
			t.FailNow()
			return
		}

		res, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, res, 0)
	})
}
