package public_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/stretchr/testify/assert"
)

func Test_Server_UserInBotRegistrationHandler(t *testing.T) {
	s, _, teardown := server.TestServer(t)
	defer teardown()

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
				"chat_id":  "",
				"username": "test",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid",
			payload: map[string]interface{}{
				"chat_id":  "invalid",
				"username": "I0HuKc",
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
				"username": "I0HuKc",
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

func Test_Server_UserAdminHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis.Registration, redis.Auth)

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

	/* Регистрация -> ШАГ 1 */

	tcRegistrationStep1 := []struct {
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
			name: "empty username",
			payload: map[string]interface{}{
				"username": "",
				"password": "4tfgefhey75uh",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "outsider username",
			payload: map[string]interface{}{
				"username": "outsider",
				"password": "4tfgefhey75uh",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty password",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"password": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "short password",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"password": "1235",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "to long password",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"password": "1235678901234567890",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"password": "4tfgefhey75uh",
				"testing":  true,
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tcRegistrationStep1 {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/registration/code", b)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	/* Регистрация -> ШАГ 2 */

	tcRegistrationStep2 := []struct {
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
			name: "empty code",
			payload: map[string]interface{}{
				"code": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid code (short)",
			payload: map[string]interface{}{
				"code": 134,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid code (long)",
			payload: map[string]interface{}{
				"code": 9764115,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"code": 100000,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range tcRegistrationStep2 {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/registration", b)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	/* Авторизация */

	tcAuth := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "empty username",
			payload: map[string]interface{}{
				"username": "",
				"password": "4tfgefhey75uh",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid username",
			payload: map[string]interface{}{
				"username": "invalid",
				"password": "4tfgefhey75uh",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty password",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"password": "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid password",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"password": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"username": "I0HuKc",
				"password": "4tfgefhey75uh",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tcAuth {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/auth", b)
			s.Router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

		})
	}
}
