package sqlstorage

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	pgx4 "github.com/jackc/pgx/v4"
	sqlstorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	yaml3 "gopkg.in/yaml.v3"
)

const configFile = "configs/config.yaml"

func TestStorage(t *testing.T) {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		t.Skip(configFile + " file does not exists")
	}

	configContent, _ := os.ReadFile(configFile)

	var config struct {
		Storage struct {
			URL string
		}
	}

	err := yaml3.Unmarshal(configContent, &config)
	if err != nil {
		t.Fatal("Failed to unmarshal config", err)
	}

	ctx := context.Background()
	storage := New(ctx, config.Storage.URL)
	if err := storage.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to DB server", err)
	}

	t.Run("test SQL", func(t *testing.T) {
		tx, err := storage.conn.BeginTx(ctx, pgx4.TxOptions{
			IsoLevel:       pgx4.Serializable,
			AccessMode:     pgx4.ReadWrite,
			DeferrableMode: pgx4.NotDeferrable,
		})
		if err != nil {
			t.Fatal("Failed to connect to DB server", err)
		}

		userID := uuid.New()
		started, err := time.Parse("2006-01-02 15:04:05", "2022-03-10 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}
		ended, err := time.Parse("2006-01-02 15:04:05", "2022-03-11 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}

		notify, err := time.Parse("2006-01-02 15:04:05", "2022-03-09 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}

		event := sqlstorage.NewEvent(
			"Title",
			started,
			ended,
			"Description",
			userID,
			notify,
		)

		err = storage.Create(*event)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err := storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		event.Title = "Test title"

		err = storage.Update(*event)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		err = storage.Delete(event.ID)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 0)

		find, err := storage.Find(event.ID)
		if err != nil {
			t.Fatal()
			return
		}
		require.Len(t, find, 1)

		err = tx.Rollback(ctx)
		if err != nil {
			t.Fatal("Failed to rollback tx", err)
		}
	})
}
