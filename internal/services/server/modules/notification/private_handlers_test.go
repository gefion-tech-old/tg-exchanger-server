package notification_test

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

func Test_Server_DeleteNotification(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	assert.NoError(t, server.TestNotification854(t, s, tokens))

	testCases := []struct {
		name         string
		id           string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "undefined id",
			id:           "",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid id",
			id:           "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid",
			id:           "1",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/notification/"+tc.id, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_UpdateNotificationStatus(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	assert.NoError(t, server.TestNotification854(t, s, tokens))

	testCases := []struct {
		name         string
		id           int
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty type",
			id:   1,
			payload: map[string]interface{}{
				"status": 2,
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid type 1",
			id:   1,
			payload: map[string]interface{}{
				"status": 2,
				"type":   "invalid",
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid type 2",
			id:   1,
			payload: map[string]interface{}{
				"status": 2,
				"type":   100,
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty status",
			id:   1,
			payload: map[string]interface{}{
				"type": 854,
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid status 1",
			id:   1,
			payload: map[string]interface{}{
				"status": "invalid",
				"type":   854,
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid status 2",
			id:   1,
			payload: map[string]interface{}{
				"status": 100,
				"type":   854,
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid user",
			id:   1,
			payload: map[string]interface{}{
				"status": 2,
				"type":   854,
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			id:   1,
			payload: map[string]interface{}{
				"status": 2,
				"type":   854,
				"meta_data": map[string]interface{}{
					"code":      245335,
					"user_card": "5559494130410854",
					"img_path":  "tmp/some_path.png",
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
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
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/admin/notification/%d", tc.id), b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_CreateNotification(t *testing.T) {
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
			name: "empty type",
			payload: map[string]interface{}{
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245335,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid type 1",
			payload: map[string]interface{}{
				"type": "invalid",
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245335,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid type 2",
			payload: map[string]interface{}{
				"type": 100,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245335,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid user",
			payload: map[string]interface{}{
				"type": 854,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245335,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid",
			payload:      mocks.ADMIN_NOTIFICATION_854,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/notification", b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_GetAllNotifications(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю уведомления
	for i := 0; i < 14; i++ {
		assert.NoError(t, server.TestNotification854(t, s, tokens))
	}

	testCases := []struct {
		name              string
		page              int
		limit             int
		expectedArrLength int
	}{
		{
			name:              "page 1",
			page:              1,
			limit:             5,
			expectedArrLength: 5,
		},
		{
			name:              "page 2",
			page:              2,
			limit:             5,
			expectedArrLength: 5,
		},
		{
			name:              "page 3",
			page:              3,
			limit:             5,
			expectedArrLength: 4,
		},
		{
			name:              "page 1",
			page:              1,
			limit:             15,
			expectedArrLength: 14,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/admin/notifications?page=%d&limit=%d", tc.page, tc.limit)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			body := map[string]interface{}{}
			decodeErr := json.NewDecoder(rec.Body).Decode(&body)
			assert.NoError(t, decodeErr)
			assert.NotNil(t, body["data"])

			assert.Equal(t, tc.expectedArrLength, len(body["data"].([]interface{})))
		})
	}
}
