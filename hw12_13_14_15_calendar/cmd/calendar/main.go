package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalapp "github.com/socialdistance/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/socialdistance/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/socialdistance/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/socialdistance/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/socialdistance/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := internalconfig.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed load config %s", err)
	}

	logg, err := internallogger.New(config.Logger)
	if err != nil {
		log.Fatalf("Failed logger %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	store := CreateStorage(ctx, *config)
	calendar := internalapp.New(logg, store)

	grpc := internalgrpc.NewServer(logg, calendar, config.GRPC.Host, config.GRPC.Port)

	go func() {
		if err := grpc.Start(); err != nil {
			logg.Error("Failed to start GRPC server: %s", err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		grpc.Stop()
	}()

	server := internalhttp.NewServer(logg, calendar, config.HTTP.Host, config.HTTP.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

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
