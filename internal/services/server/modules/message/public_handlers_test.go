package message_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/stretchr/testify/assert"
)

func Test_Server_GetMessageHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю тестовое сообщение
	assert.NoError(t, server.TestBotMessage(t, s, tokens))

	testCases := []struct {
		name         string
		connector    string
		expectedCode int
	}{
		{
			name:         "undefined connector",
			connector:    "undefined",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "valid",
			connector:    mocks.BOT_MESSAGE_REQ["connector"].(string),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/message/"+tc.connector, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusOK {
				t.Run("response_validation", func(t *testing.T) {
					var body models.BotMessage
					assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
					assert.NotNil(t, body)
					assert.NoError(t, body.Validation())
				})
			}
		})
	}
}
