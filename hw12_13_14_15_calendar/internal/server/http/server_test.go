package internalhttp

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestHttpServerHelloWorld(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	httpHandlers := NewRouter(createApp(t))
	httpHandlers.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, "Hello, world!\n", string(body))
}

func TestHttpServerEventsOperations(t *testing.T) {
	body := bytes.NewBufferString(`{
		"id": "dd983962-b469-11ec-b909-0242ac120002",
		"userId": "fae10b5c-b469-11ec-b909-0242ac120002",
		"title": "Test Title",
		"startedAt": "2022-04-01 00:00:00",
		"finishedAt": "2022-05-01 00:00:00",
		"description": "Test Description",
		"notifyAt": "2022-03-30 00:00:00"
	}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	httpHandlers := NewRouter(createApp(t))
	httpHandlers.ServeHTTP(w, req)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)
	respExpected := `{"id":"dd983962-b469-11ec-b909-0242ac120002","userId":"fae10b5c-b469-11ec-b909-0242ac120002","title":"Test Title","startedAt":"2022-04-01 00:00:00","finishedAt":"2022-05-01 00:00:00","description":"Test Description","notifyAt":"2022-03-30 00:00:00"}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))

	body = bytes.NewBufferString(`{
		"userId": "e648562c-b46a-11ec-b909-0242ac120002",
		"title": "New Test Title",
		"startedAt": "2023-04-01 00:00:00",
		"finishedAt": "2023-05-01 00:00:00",
		"description": "New Test Description",	
		"notifyAt": "2023-03-30 00:00:00"
	}`)
	req = httptest.NewRequest("PUT", "/events/dd983962-b469-11ec-b909-0242ac120002", body)
	w = httptest.NewRecorder()

	httpHandlers.ServeHTTP(w, req)

	resp = w.Result()
	respBody, _ = io.ReadAll(resp.Body)
	respExpected = `{"id":"dd983962-b469-11ec-b909-0242ac120002","userId":"e648562c-b46a-11ec-b909-0242ac120002","title":"New Test Title","startedAt":"2023-04-01 00:00:00","finishedAt":"2023-05-01 00:00:00","description":"New Test Description","notifyAt":"2023-03-30 00:00:00"}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))
}

func createApp(t *testing.T) *app.App {
	t.Helper()
	logFile, err := os.CreateTemp("", "test.log")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	logger, err := logger.New(internalconfig.LoggerConf{Level: "info", File: logFile.Name()})
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	inMemoryStorage := memorystorage.New()

	return app.New(logger, inMemoryStorage)
}
