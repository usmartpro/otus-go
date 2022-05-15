package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/server/http"
	internalfactory "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	configuration, err := internalconfig.LoadCalendarConfiguration(configFile)
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}
	logg, err := internallogger.New(configuration.Logger)
	if err != nil {
		logg.Error("error create logger: " + err.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := internalfactory.NewStorage(ctx, configuration.Storage)
	calendar := app.New(logg, storage)

	// gRPC
	serverGrpc := internalgrpc.NewServer(logg, calendar, configuration.GRPC.Host, configuration.GRPC.Port)

	go func() {
		if err := serverGrpc.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		serverGrpc.Stop()
	}()

	// HTTP
	server := internalhttp.NewServer(logg, calendar, configuration.HTTP.Host, configuration.HTTP.Port)

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	<-ctx.Done()
}
