package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(message string, params ...interface{})
	Info(message string, params ...interface{})
	Error(message string, params ...interface{})
	Warn(message string, params ...interface{})
}

type Storage interface {
	Create(e storage.Event) error
	Update(e storage.Event) error
	Delete(id uuid.UUID) error
	FindAll() ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, evt storage.Event) error {
	a.Logger.Info("Creating new event")

	if err := a.Storage.Create(evt); err != nil {
		a.Logger.Error("Create new event error: %s", err)
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, evt storage.Event) error {
	a.Logger.Info("Updating event")

	if err := a.Storage.Update(evt); err != nil {
		a.Logger.Error("Update event error: %s", err)
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	a.Logger.Info("Deleting event with id:%s", id)

	if err := a.Storage.Delete(id); err != nil {
		a.Logger.Error("Delete event error: %s", err)
		return err
	}

	return nil
}

func (a *App) FindAllEvent(ctx context.Context) ([]storage.Event, error) {
	return a.Storage.FindAll()
}
