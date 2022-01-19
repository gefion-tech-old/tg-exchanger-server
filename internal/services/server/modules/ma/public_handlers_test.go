package ma_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/stretchr/testify/assert"
)

func Test_Server_GetMerchantAutopayoutSelectionHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	assert.NoError(t, server.TestMerchantAutopayout(t, s, tokens))

	testCases := []struct {
		name         string
		page         string
		limit        string
		service      string
		expectedCode int
	}{
		{
			name:         "invalid page",
			page:         "invalid",
			limit:        "15",
			service:      AppType.MerchantAutoPayoutWhitebit,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid page",
			page:         "0",
			limit:        "15",
			service:      AppType.MerchantAutoPayoutWhitebit,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid limit",
			page:         "1",
			limit:        "invalid",
			service:      AppType.MerchantAutoPayoutWhitebit,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid limit",
			page:         "1",
			limit:        "1000",
			service:      AppType.MerchantAutoPayoutWhitebit,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid service",
			page:         "1",
			limit:        "10",
			service:      "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid",
			page:         "1",
			limit:        "15",
			service:      AppType.MerchantAutoPayoutWhitebit,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet,
				fmt.Sprintf("/api/v1/admin/merchant-autopayout/all?page=%s&limit=%s&service=%s", tc.page, tc.limit, tc.service), nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func Test_Server_DeleteMerchantAutopayoutHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	assert.NoError(t, server.TestMerchantAutopayout(t, s, tokens))

	testCases := []struct {
		name         string
		id           string
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
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/merchant-autopayout/"+tc.id, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusOK {
				var body models.MerchantAutopayout
				assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

				t.Run("response_validation", func(t *testing.T) {
					assert.NotNil(t, body)
					assert.NoError(t, body.Validation())
				})
			}
		})
	}
}

func Test_Server_UpdateMerchantAutopayoutHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	assert.NoError(t, server.TestMerchantAutopayout(t, s, tokens))

	testCases := []struct {
		name         string
		id           string
		payload      interface{}
		expectedCode int
	}{
		{
			name:         "undefined id",
			id:           "",
			payload:      mocks.MerchantAutopayout,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid id",
			id:           "invalid",
			payload:      mocks.MerchantAutopayout,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid payload",
			id:           "1",
			payload:      "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid",
			id:           "1",
			payload:      mocks.MerchantAutopayout,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/merchant-autopayout/"+tc.id, b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusOK {
				var body models.MerchantAutopayout
				assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

				t.Run("response_validation", func(t *testing.T) {
					assert.NotNil(t, body)
					assert.NoError(t, body.Validation())
				})
			}
		})
	}
}

func Test_Server_CreateMerchantAutopayoutHandler(t *testing.T) {
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
			name: "empty name",
			payload: map[string]interface{}{
				// "name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid name",
			payload: map[string]interface{}{
				"name":         "main accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty service",
			payload: map[string]interface{}{
				"name": "main_accout",
				// "service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid service",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      "invalid",
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty service_type",
			payload: map[string]interface{}{
				"name":    "main_accout",
				"service": AppType.MerchantAutoPayoutWhitebit,
				// "service_type": AppType.UseAsMerchant,
				"options":    "{}",
				"status":     true,
				"message_id": 1,
				"created_by": "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid service_type",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": "invalid",
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty options",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				// "options":      "{}",
				"status":     true,
				"message_id": 1,
				"created_by": "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid options",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "invalid",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty status",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				// "status":       true,
				"message_id": 1,
				"created_by": "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid status",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       "invalid",
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty message_id",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				// "message_id":   1,
				"created_by": "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid message_id",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   "invalid",
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty created_by",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				// "created_by":   "I0HuKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid created_by",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0H uKc",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid",
			payload: map[string]interface{}{
				"name":         "main_accout",
				"service":      AppType.MerchantAutoPayoutWhitebit,
				"service_type": AppType.UseAsMerchant,
				"options":      "{}",
				"status":       true,
				"message_id":   1,
				"created_by":   "I0HuKc",
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/merchant-autopayout", b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusCreated {
				var body models.MerchantAutopayout
				assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

				t.Run("response_validation", func(t *testing.T) {
					assert.NotNil(t, body)
					assert.NoError(t, body.Validation())
				})
			}
		})
	}
}
