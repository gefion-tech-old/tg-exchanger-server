package exchanger_test

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

func Test_Server_DeleteExchangerHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	assert.NoError(t, server.TestExchanger(t, s, tokens))

	testCases := []struct {
		name         string
		id           string
		expectedCode int
	}{
		{
			name:         "undefined id",
			id:           "10",
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
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/exchanger/"+tc.id, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_UpdateExchangerHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	assert.NoError(t, server.TestExchanger(t, s, tokens))

	testCases := []struct {
		name         string
		id           string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "invalid payload",
			id:           "1",
			payload:      "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty name",
			id:   "1",
			payload: map[string]interface{}{
				"name": "",
				"url":  "http://some.com",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid name",
			id:   "1",
			payload: map[string]interface{}{
				"name": "1 obmen",
				"url":  "http://some.com",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty url",
			id:   "1",
			payload: map[string]interface{}{
				"name": "1obmen",
				"url":  "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid url",
			id:   "1",
			payload: map[string]interface{}{
				"name": "1obmen",
				"url":  "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "undefined id",
			id:   "10",
			payload: map[string]interface{}{
				"name": "1obmen",
				"url":  "http://some.com",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "invalid id",
			id:   "invalid",
			payload: map[string]interface{}{
				"name": "1obmen",
				"url":  "http://some.com",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid",
			id:           "1",
			payload:      mocks.ADMIN_EXCHANGER,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/exchanger/"+tc.id, b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_CreateExchangerHandler(t *testing.T) {
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
			name: "empty created_by",
			payload: map[string]interface{}{
				"created_by": "",
				"name":       "1obmen",
				"url":        "https://1obmen.net/request-exportxml.xml",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty name",
			payload: map[string]interface{}{
				"created_by": "I0HuKc",
				"name":       "",
				"url":        "https://1obmen.net/request-exportxml.xml",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "short name",
			payload: map[string]interface{}{
				"created_by": "I0HuKc",
				"name":       "12",
				"url":        "https://1obmen.net/request-exportxml.xml",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid name",
			payload: map[string]interface{}{
				"created_by": "I0HuKc",
				"name":       "1 23",
				"url":        "https://1obmen.net/request-exportxml.xml",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty url",
			payload: map[string]interface{}{
				"created_by": "I0HuKc",
				"name":       "1obmen",
				"url":        "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid url",
			payload: map[string]interface{}{
				"name": "1obmen",
				"url":  "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid",
			payload:      mocks.ADMIN_EXCHANGER,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Кодирую тело запроса
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/exchanger", b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_GetExchangersSelectionHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю уведомления
	for i := 0; i < 14; i++ {
		assert.NoError(t, server.TestExchanger(t, s, tokens))
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
			url := fmt.Sprintf("/api/v1/admin/exchangers?page=%d&limit=%d", tc.page, tc.limit)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			body := map[string]interface{}{}
			decodeErr := json.NewDecoder(rec.Body).Decode(&body)
			assert.NoError(t, decodeErr)
			assert.NotNil(t, body["data"])

			assert.Equal(t, len(body["data"].([]interface{})), tc.expectedArrLength)
		})
	}
}
