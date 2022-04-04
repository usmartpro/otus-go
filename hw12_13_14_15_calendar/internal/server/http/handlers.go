package internalhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/usmartpro/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type ServerHandlers struct {
	app *app.App
}

func NewServerHandlers(a *app.App) *ServerHandlers {
	return &ServerHandlers{app: a}
}

func (s *ServerHandlers) HelloWorld(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello, world!\n")
	w.WriteHeader(200)
	w.Write(msg)
}

func (s *ServerHandlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var dto Dto
	err := ParseRequest(r, &dto)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	event, err := dto.GetModel()
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.CreateEvent(r.Context(), *event)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	responseData, _ := json.Marshal(dto)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func (s *ServerHandlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var dto Dto
	err := ParseRequest(r, &dto)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	dto.ID = vars["id"]

	event, err := dto.GetModel()
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.UpdateEvent(r.Context(), *event)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	responseData, _ := json.Marshal(dto)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func (s *ServerHandlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
	}

	err = s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHandlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	interval := r.URL.Query().Get("interval")

	withDate := false
	dtStart, err := time.Parse("2006-01-02", date)
	if err == nil {
		withDate = true
	}

	var events []storage.Event
	if withDate {
		switch interval {
		case "day":
			events, err = s.app.GetDayEvents(r.Context(), dtStart)
		case "week":
			events, err = s.app.GetWeekEvents(r.Context(), dtStart)
		case "month":
			events, err = s.app.GetMonthEvents(r.Context(), dtStart)
		default:
			events, err = s.app.GetDayEvents(r.Context(), dtStart)
		}
	} else {
		events, err = s.app.GetEvents(r.Context())
	}

	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
	}

	eventDtos := make([]Dto, 0, len(events))
	for _, e := range events {
		eventDtos = append(eventDtos, ModelToEventDto(e))
	}

	responseData, _ := json.Marshal(eventDtos)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func ParseRequest(r *http.Request, dto interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(data, dto)
	if err != nil {
		return fmt.Errorf("failed to decode JSON request: %w", err)
	}

	return nil
}

func RespondError(w http.ResponseWriter, code int, err error) {
	data, err := json.Marshal(ErrorDto{
		false,
		err.Error(),
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to marshall error dto"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
