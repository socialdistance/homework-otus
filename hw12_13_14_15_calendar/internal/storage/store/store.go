package store

import (
	"context"
	"log"

	internalapp "github.com/socialdistance/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/socialdistance/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage/sql"
)

func CreateStorage(ctx context.Context, config internalconfig.Config) internalapp.Storage {
	var store internalapp.Storage

	switch config.Storage.Type {
	case internalconfig.InMemmory:
		store = memorystorage.New()
	case internalconfig.SQL:
		sqlStore := sqlstorage.New(ctx, config.Storage.URL)
		err := sqlStore.Connect(ctx)
		if err != nil {
			log.Fatalf("Unable to connect database: %s", err)
		}
		store = sqlStore
	default:
		log.Fatalf("Dont know type storage: %s", config.Storage.Type)
	}

	return store
}
