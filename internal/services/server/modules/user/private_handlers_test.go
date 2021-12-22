package user_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/stretchr/testify/assert"
)

func Test_Server_LogoutHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "empty token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid token",
			token:        fmt.Sprintf("Bearer %s", "invalid"),
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "valid",
			token:        fmt.Sprintf("Bearer %s", tokens["access_token"]),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/logout", nil)
			req.Header.Add("Authorization", tc.token)
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
