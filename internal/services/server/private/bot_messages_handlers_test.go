package private_test

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

func Test_Server_DeleteBotMessageHandler(t *testing.T) {
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
			name:         "empty connector",
			connector:    "",
			expectedCode: http.StatusNotFound,
		},
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
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/message/"+tc.connector, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}

func Test_Server_UpdateBotMessageHandler(t *testing.T) {
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
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "invalid payload",
			connector:    mocks.BOT_MESSAGE_REQ["connector"].(string),
			payload:      "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "empty connector",
			connector: "",
			payload: map[string]interface{}{
				"message_text": mocks.BOT_MESSAGE_REQ["message_text"],
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:      "empty message_text",
			connector: mocks.BOT_MESSAGE_REQ["connector"].(string),
			payload: map[string]interface{}{
				"message_text": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "valid",
			connector: mocks.BOT_MESSAGE_REQ["connector"].(string),
			payload: map[string]interface{}{
				"message_text": "new text",
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
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/message/"+tc.connector, b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusOK {
				// Проверяю обновился ли текст
				body := models.BotMessage{}
				decodeErr := json.NewDecoder(rec.Body).Decode(&body)
				assert.NoError(t, decodeErr)
				assert.Equal(t, "new text", body.MessageText)
			}
		})
	}

}

func Test_Server_GetAllBotMessageHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю тестовое сообщение
	assert.NoError(t, server.TestBotMessage(t, s, tokens))

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/messages", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
	s.Router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body := map[string]interface{}{}
	decodeErr := json.NewDecoder(rec.Body).Decode(&body)
	assert.NoError(t, decodeErr)
	assert.Len(t, body["data"], 1)

}

func Test_Server_GetBotMessageHandler(t *testing.T) {
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
		})
	}

}

func Test_Server_CreateNewBotMessageHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

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
			name: "empty connector",
			payload: map[string]interface{}{
				"connector":    "",
				"message_text": mocks.BOT_MESSAGE_REQ["message_text"],
				"created_by":   mocks.BOT_MESSAGE_REQ["created_by"],
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid connector 1",
			payload: map[string]interface{}{
				"connector":    "one two",
				"message_text": mocks.BOT_MESSAGE_REQ["message_text"],
				"created_by":   mocks.BOT_MESSAGE_REQ["created_by"],
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid connector 2",
			payload: map[string]interface{}{
				"connector":    "one..two",
				"message_text": mocks.BOT_MESSAGE_REQ["message_text"],
				"created_by":   mocks.BOT_MESSAGE_REQ["created_by"],
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty message_text",
			payload: map[string]interface{}{
				"connector":    mocks.BOT_MESSAGE_REQ["connector"],
				"message_text": "",
				"created_by":   mocks.BOT_MESSAGE_REQ["created_by"],
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty created_by",
			payload: map[string]interface{}{
				"connector":    mocks.BOT_MESSAGE_REQ["connector"],
				"message_text": mocks.BOT_MESSAGE_REQ["message_text"],
				"created_by":   "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"connector":    mocks.BOT_MESSAGE_REQ["connector"],
				"message_text": mocks.BOT_MESSAGE_REQ["message_text"],
				"created_by":   mocks.BOT_MESSAGE_REQ["created_by"],
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/message", b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
