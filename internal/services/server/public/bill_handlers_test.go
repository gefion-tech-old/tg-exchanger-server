package public_test

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

func Test_Server_GetAllBillsHandler(t *testing.T) {
	s, _, teardown := server.TestServer(t)
	defer teardown()

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

	// Создаю тестовый пользовательский счет
	assert.NoError(t, server.TestUserBill(t, s))

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
	s, _, teardown := server.TestServer(t)
	defer teardown()

	// Регистрирую пользователя в боте
	assert.NoError(t, server.TestBotUser(t, s))

	// Создаю тестовый пользовательский счет
	assert.NoError(t, server.TestUserBill(t, s))

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

func Test_Server_NewBillHandler(t *testing.T) {
	s, _, teardown := server.TestServer(t)
	defer teardown()

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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/bot/user/bill", b)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
