package logs_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/stretchr/testify/assert"
)

func Test_Server_CreateLogRecordHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty service",
			payload: map[string]interface{}{
				// "service": 235,
				"module": "db",
				"info":   "some error text",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty module for server log",
			payload: map[string]interface{}{
				"service": static.L__SERVER,
				// "module": "db",
				"info": "some error text",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty info for server log",
			payload: map[string]interface{}{
				"service": static.L__SERVER,
				"module":  "db",
				// "info":    "some error text",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},

		{
			name: "empty username for admin log",
			payload: map[string]interface{}{
				"service": static.L__ADMIN,
				"module":  "db",
				"info":    "some error text",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid service",
			payload: map[string]interface{}{
				"service": 235,
				"module":  "db",
				"info":    "some error text",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid server log",
			payload: map[string]interface{}{
				"service": static.L__SERVER,
				"module":  "db",
				"info":    "some error text",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "valid bot log",
			payload: map[string]interface{}{
				"service": static.L__BOT,
				"module":  "db",
				"info":    "some error text",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "valid admin log",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"service":  static.L__ADMIN,
				"module":   "db",
				"info":     "some error text",
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/log", b)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
