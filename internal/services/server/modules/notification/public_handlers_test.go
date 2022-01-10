package notification_test

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
				"status": 1,
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
				"type":   "invalid",
				"status": 1,
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
				"type":   100,
				"status": 1,
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
				"type":   854,
				"status": 1,
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
			name: "invalid verification code",
			payload: map[string]interface{}{
				"type":   854,
				"status": 1,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      24533,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid user card",
			payload: map[string]interface{}{
				"type":   854,
				"status": 1,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245332,
						"user_card": "5559494130410",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid username",
			payload: map[string]interface{}{
				"type":   854,
				"status": 1,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245331,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": map[string]interface{}{
					"chat_id":  3673563,
					"username": "<script> </script>",
				},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid chat_id",
			payload: map[string]interface{}{
				"type":   854,
				"status": 1,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245331,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": map[string]interface{}{
					"chat_id":  "invalid",
					"username": "I0HuKc",
				},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty chat_id",
			payload: map[string]interface{}{
				"type":   854,
				"status": 1,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245331,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": map[string]interface{}{
					// "chat_id":  "invalid",
					"username": "I0HuKc",
				},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty username",
			payload: map[string]interface{}{
				"type":   854,
				"status": 1,
				"meta_data": map[string]interface{}{
					"card_verification": map[string]interface{}{
						"code":      245331,
						"user_card": "5559494130410854",
						"img_path":  "tmp/some_path.png",
					},
				},
				"user": map[string]interface{}{
					"chat_id": 3673563,
					// "username": "<script></script>",
				},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid",
			payload:      mocks.ADMIN_NOTIFICATION_854,
			expectedCode: http.StatusCreated,
		},

		{
			name: "invalid meta_data",
			payload: map[string]interface{}{
				"type":      855,
				"status":    1,
				"meta_data": "invalid",
				"user":      mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty ex_from",
			payload: map[string]interface{}{
				"type":   855,
				"status": 1,
				"meta_data": map[string]interface{}{
					"ex_action_cancel": map[string]interface{}{
						// "ex_from": "ETH",
						"ex_to": "BTC",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty ex_to",
			payload: map[string]interface{}{
				"type":   855,
				"status": 1,
				"meta_data": map[string]interface{}{
					"ex_action_cancel": map[string]interface{}{
						"ex_from": "ETH",
						// "ex_to": "BTC",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "valid action cancel",
			payload: map[string]interface{}{
				"type":   855,
				"status": 1,
				"meta_data": map[string]interface{}{
					"ex_action_cancel": map[string]interface{}{
						"ex_from": "ETH",
						"ex_to":   "BTC",
					},
				},
				"user": mocks.USER_IN_BOT_REGISTRATION_REQ,
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/notification", b)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusCreated {
				t.Run("response_validation", func(t *testing.T) {
					var body models.Notification
					assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
					assert.NotNil(t, body)
					assert.NoError(t, body.Validation())
				})
			}
		})
	}
}
