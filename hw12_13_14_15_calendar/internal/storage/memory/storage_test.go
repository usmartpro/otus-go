package memorystorage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	memorystorage "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	storage := New()

	t.Run("crud test", func(t *testing.T) {
		userID := uuid.New()
		startedAt, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 00:00:00")
		finishedAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
		notifyBeforeTime, _ := time.Parse("2006-01-02 15:04:05", "2022-06-01 00:00:00")

		event := memorystorage.NewEvent(
			userID,
			"Test name",
			startedAt,
			finishedAt,
			"Test description",
			notifyBeforeTime,
		)

		// Insert
		_ = storage.Insert(*event)

		bdRecord, err := storage.Select()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, bdRecord, 1)
		require.Equal(t, *event, bdRecord[0])

		// Update
		event.Title = "Test name 2"
		event.Description = ""

		err = storage.Update(*event)
		if err != nil {
			t.FailNow()
			return
		}

		bdRecord, err = storage.Select()
		if err != nil {
			t.FailNow()
			return
		}

		require.Len(t, bdRecord, 1)
		require.Equal(t, *event, bdRecord[0])
		require.Equal(t, event.Title, bdRecord[0].Title)
		require.Equal(t, event.Description, bdRecord[0].Description)

		// Delete
		_ = storage.Delete(event.ID)
		bdRecord, err = storage.Select()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, bdRecord, 0)
	})
}
