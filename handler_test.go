package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndListTasks(t *testing.T) {
	storage := NewStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", storage.CreateHandler)
	mux.HandleFunc("GET /tasks", storage.ListHandler)

	// 1. Создаём задачу
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"text":"тест"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Получаем список
	req = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var tasks []Task
	json.NewDecoder(w.Body).Decode(&tasks)

	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "тест", tasks[0].Text)
}

func TestCreateHandler_TableDriven(t *testing.T) {
	storage := NewStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", storage.CreateHandler)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedTasks  int // сколько задач должно стать после запроса
	}{
		{
			name:           "успешное создание",
			payload:        `{"text":"купить молоко"}`,
			expectedStatus: http.StatusCreated,
			expectedTasks:  1,
		},
		{
			name:           "пустой текст — ошибка",
			payload:        `{"text":""}`,
			expectedStatus: http.StatusBadRequest,
			expectedTasks:  0,
		},
		{
			name:           "невалидный JSON",
			payload:        `{"text":}`,
			expectedStatus: http.StatusBadRequest,
			expectedTasks:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// каждый кейс — в отдельном подтесте (видно в go test -v)
			storage := NewStorage() // чистое хранилище для каждого кейса
			mux := http.NewServeMux()
			mux.HandleFunc("POST /tasks", storage.CreateHandler)

			req := httptest.NewRequest("POST", "/tasks", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedTasks, len(storage.List()))
		})
	}
}
