package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/mq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/sender_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	configuration, err := internalconfig.LoadSenderConfig(configFile)
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}

	logg, err := internallogger.New(configuration.Logger)
	if err != nil {
		log.Fatalf("error create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

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

	sender := app.NewSender(rabbitClient, logg)
	sender.Run()

	<-ctx.Done()
}
