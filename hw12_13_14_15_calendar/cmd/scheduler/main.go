package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/mq"
	internalfactory "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	configuration, err := internalconfig.LoadSchedulerConfig(configFile)
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}

	logg, err := internallogger.New(configuration.Logger)
	if err != nil {
		log.Fatalf("error create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := internalfactory.NewStorage(ctx, configuration.Storage)

	rabbitClient, err := mq.NewRabbit(
		ctx,
		configuration.Rabbit.Dsn,
		configuration.Rabbit.Exchange,
		configuration.Rabbit.Queue,
		logg)
	if err != nil {
		cancel()
		log.Fatalf("error create rabbit client: %s", err) //nolint:gocritic
	}

	scheduler := app.NewAppScheduler(storage.(app.EventSource), rabbitClient, logg)

	timer := time.Tick(time.Second)
	timerHour := time.Tick(time.Hour)

	go func() {
		for {
			select {
			case <-timer:
				err := scheduler.Notify()
				if err != nil {
					logg.Error("error Notify: %s", err)
				}
			case <-timerHour:
				err := scheduler.RemoveOldEvents()
				if err != nil {
					logg.Error("error RemoveOldEvents: %s", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	logg.Info("scheduler is running...")

	<-ctx.Done()
}
