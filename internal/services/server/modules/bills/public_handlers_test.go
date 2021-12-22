package bills_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
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
				"chat_id": mocks.USER_BILL_REQ["chat_id"],
				"bill":    "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid bill format",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_BILL_REQ["chat_id"],
				"bill":    "5559493 130410854",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid bill lenght (short)",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_BILL_REQ["chat_id"],
				"bill":    "55594931304108",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid bill lenght (long)",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_BILL_REQ["chat_id"],
				"bill":    "5559493130410854000",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "not found bill",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_BILL_REQ["chat_id"],
				"bill":    "5559494130410829",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"chat_id": mocks.USER_BILL_REQ["chat_id"],
				"bill":    mocks.USER_BILL_REQ["bill"],
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
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/bot/user/bill", b)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
