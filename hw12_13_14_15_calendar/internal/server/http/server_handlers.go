package internalhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	internalapp "github.com/socialdistance/hw12_13_14_15_calendar/internal/app"
)

type ServerHandlers struct {
	app *internalapp.App
}

func ResponseError(w http.ResponseWriter, code int, err error) {
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

func NewServerHandlers(a *internalapp.App) *ServerHandlers {
	return &ServerHandlers{app: a}
}

func (s *ServerHandlers) HelloWorld(w http.ResponseWriter, r *http.Request) {
	response := []byte("Hello world!\n")
	w.WriteHeader(200)
	w.Write(response)
}

func (s *ServerHandlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var dto EventDto

	err := ParsingData(r, &dto)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	event, err := dto.GetModel()
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.CreateEvent(r.Context(), *event)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(dto)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (s *ServerHandlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var dto EventDto
	err := ParsingData(r, &dto)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	dto.ID = vars["id"]

	event, err := dto.GetModel()
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.UpdateEvent(r.Context(), *event)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(dto)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (s *ServerHandlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		ResponseError(w, http.StatusBadGateway, err)
		return
	}

	err = s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		ResponseError(w, http.StatusBadGateway, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHandlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := s.app.FindAllEvent(r.Context())
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}
	eventDto := make([]EventDto, 0, len(events))
	for _, t := range events {
		eventDto = append(eventDto, CreateDto(t))
	}

	response, err := json.Marshal(eventDto)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ParsingData(r *http.Request, dto interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error read body: %w", err)
	}

	err = json.Unmarshal(data, dto)
	if err != nil {
		return fmt.Errorf("error unmarshall body: %w", err)
	}

	return nil
}
