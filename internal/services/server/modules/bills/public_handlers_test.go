package bills_test

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

func Test_Server_GetBillHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

	// Создаю тестовый пользовательский счет
	assert.NoError(t, server.TestUserBill(t, s, nil))

	testCases := []struct {
		name         string
		id           string
		expectedCode int
	}{
		{
			name:         "invalid id",
			id:           "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "undefined id",
			id:           "10",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/bot/user/bill/"+tc.id, nil)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_GetAllBillsHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

	// Создаю тестовый пользовательский счет
	assert.NoError(t, server.TestUserBill(t, s, nil))

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "valid",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/bot/user/%d/bills", mocks.USER_BILL_REQ["chat_id"]), nil)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}

func Test_Server_DeleteBillHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю тестовый пользовательский счет
	assert.NoError(t, server.TestUserBill(t, s, tokens))

	testCases := []struct {
		name    string
		payload struct {
			chat_id string
			bill_id string
		}
		expectedCode int
	}{
		{
			name: "empty chat_id",
			payload: struct {
				chat_id string
				bill_id string
			}{
				bill_id: "1",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid chat_id",
			payload: struct {
				chat_id string
				bill_id string
			}{
				chat_id: "invalid",
				bill_id: "1",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty bill_id",
			payload: struct {
				chat_id string
				bill_id string
			}{
				chat_id: fmt.Sprintf("%d", mocks.USER_BILL_REQ["chat_id"]),
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "invalid bill_id",
			payload: struct {
				chat_id string
				bill_id string
			}{
				chat_id: fmt.Sprintf("%d", mocks.USER_BILL_REQ["chat_id"]),
				bill_id: "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},

		{
			name: "empty payload",
			payload: struct {
				chat_id string
				bill_id string
			}{},
			expectedCode: http.StatusNotFound,
		},

		{
			name: "valid",
			payload: struct {
				chat_id string
				bill_id string
			}{
				chat_id: fmt.Sprintf("%d", mocks.USER_BILL_REQ["chat_id"]),
				bill_id: "1",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/bot/user/%s/bill/%s", tc.payload.chat_id, tc.payload.bill_id), nil)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusOK {
				t.Run("response_validation", func(t *testing.T) {
					var body models.Bill
					assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
					assert.NotNil(t, body)
					assert.NoError(t, body.StructFullness())
				})
			}
		})
	}
}
