package bills_test

import (
	"bytes"
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

func Test_Server_RejectBill(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

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
			name: "empty reason",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"],
				"bill":    "5559493130410854",
				"reason":  "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"],
				"bill":    "5559493130410854",
				"reason":  "some text",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/bill/reject", b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_CreateBill(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

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
			name: "empty chat_id",
			payload: map[string]interface{}{
				"chat_id": "",
				"bill":    "5559493130410854",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty bill",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"],
				"bill":    "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid bill format",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"],
				"bill":    "5559493 130410854",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid bill lenght (short)",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"],
				"bill":    "55594931304108",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid bill lenght (long)",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"],
				"bill":    "5559493130410854000",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"],
				"bill":    "5559493130410854",
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/bill", b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusCreated {
				t.Run("response_validation", func(t *testing.T) {
					var body models.Bill
					assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
					assert.NotNil(t, body)
					assert.NoError(t, body.Validation())
				})
			}
		})
	}
}
