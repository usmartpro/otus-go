package factory

import (
	"context"
	"log"

	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage(ctx context.Context, configuration internalconfig.StorageConf) app.Storage {
	var storage app.Storage

	switch configuration.Type {
	case "memory":
		storage = memorystorage.New()
	case "base":
		storage = sqlstorage.New(ctx, configuration.Dsn).Connect(ctx)
	default:
		log.Fatalln("Unknown type of storage: " + configuration.Type)
	}

	return storage
}
