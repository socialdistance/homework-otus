package internalhttp

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/socialdistance/hw12_13_14_15_calendar/internal/app"
	"github.com/socialdistance/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/socialdistance/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/socialdistance/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	testHandlers := Routers(application(t))
	testHandlers.ServeHTTP(w, request)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	require.Equal(t, "Hello world!\n", string(body))
}

func TestHandlers(t *testing.T) {
	body := bytes.NewBufferString(`{
		"id": "a17b3f01-fbd7-40e5-8d8e-9b4cf1ef21b0",
		"title": "Test title",
		"started": "2020-10-20 12:30:00",
		"ended": "2020-10-21 12:30:00",
		"description": "test description",
		"userID": "1528371b-229c-4370-839a-0571d969902a"
	}`)

	request := httptest.NewRequest("POST", "/create", body)
	w := httptest.NewRecorder()

	testHandlers := Routers(application(t))
	testHandlers.ServeHTTP(w, request)

	response := w.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	responseExcepted := `{"id":"a17b3f01-fbd7-40e5-8d8e-9b4cf1ef21b0","title":"Test title","started":"2020-10-20 12:30:00","ended":"2020-10-21 12:30:00","description":"test description","userID":"1528371b-229c-4370-839a-0571d969902a"}` //nolint:lll
	require.Equal(t, responseExcepted, string(responseBody))

	body = bytes.NewBufferString(`{
		"title": "Test title 2",
		"started": "2020-10-20 12:30:00",
		"ended": "2020-10-21 12:30:00",
		"description": "test description 2",
		"userID": "1528371b-229c-4370-839a-0571d969902a"
	}`)

	request = httptest.NewRequest("PUT", "/events/update/a17b3f01-fbd7-40e5-8d8e-9b4cf1ef21b0", body)
	w = httptest.NewRecorder()

	testHandlers.ServeHTTP(w, request)

	response = w.Result()
	responseBody, _ = ioutil.ReadAll(response.Body)
	responseExcepted = `{"id":"a17b3f01-fbd7-40e5-8d8e-9b4cf1ef21b0","title":"Test title 2","started":"2020-10-20 12:30:00","ended":"2020-10-21 12:30:00","description":"test description 2","userID":"1528371b-229c-4370-839a-0571d969902a"}` //nolint:lll
	require.Equal(t, responseExcepted, string(responseBody))

	request = httptest.NewRequest("DELETE", "/events/delete/a17b3f01-fbd7-40e5-8d8e-9b4cf1ef21b0", body)
	w = httptest.NewRecorder()

	testHandlers.ServeHTTP(w, request)

	response = w.Result()
	responseBody, _ = ioutil.ReadAll(response.Body)
	responseExcepted = ""
	require.Equal(t, responseExcepted, string(responseBody))

	body = bytes.NewBufferString(`{
		"id": "a17b3f01-fbd7-40e5-8d8e-9b4cf1ef21b0",
		"title": "Test title",
		"started": "2020-10-20 12:30:00",
		"ended": "2020-10-21 12:30:00",
		"description": "test description",
		"userID": "1528371b-229c-4370-839a-0571d969902a"
	}`)

	requestCreate := httptest.NewRequest("POST", "/create", body)
	w = httptest.NewRecorder()

	testHandlers.ServeHTTP(w, requestCreate)

	requestListAll := httptest.NewRequest("GET", "/events", body)
	w = httptest.NewRecorder()

	testHandlers.ServeHTTP(w, requestListAll)

	response = w.Result()
	responseBody, _ = ioutil.ReadAll(response.Body)
	responseExcepted = `[{"id":"a17b3f01-fbd7-40e5-8d8e-9b4cf1ef21b0","title":"Test title","started":"2020-10-20T12:30:00Z","ended":"2020-10-21T12:30:00Z","description":"test description","userID":"1528371b-229c-4370-839a-0571d969902a"}]` // nolint:lll
	require.Equal(t, responseExcepted, string(responseBody))
}

func application(t *testing.T) *app.App {
	t.Helper()
	logFile, err := os.CreateTemp("", "log")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	logg, err := internallogger.New(config.LoggerConf{
		Level:    "info",
		Filename: logFile.Name(),
	})
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	memmoryStorage := memorystorage.New()

	return app.New(logg, memmoryStorage)
}
