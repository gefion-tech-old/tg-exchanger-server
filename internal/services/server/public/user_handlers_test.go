package public_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/go-playground/assert/v2"
)

func Test_Server_UserInBotRegistrationHandler(t *testing.T) {
	s, redis := server.TestServer(t)
	defer redis.FlushAllAsync()
	defer redis.Close()

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "empty chat_id",
			payload: map[string]interface{}{
				"chat_id":  "",
				"username": "test",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid",
			payload: map[string]interface{}{
				"chat_id":  "invalid",
				"username": "test",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty username",
			payload: map[string]interface{}{
				"chat_id":  3673563,
				"username": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"chat_id":  3673563,
				"username": "username",
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/bot/registration", b)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}
