package logs_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server"
	"github.com/stretchr/testify/assert"
)

func Test_Server_DeleteLogRecordsSelectionHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю лог записи
	for i := 0; i < 5; i++ {
		assert.NoError(t, server.TestLogRecord(t, s))
	}

	testCases := []struct {
		name         string
		from         string
		to           string
		expectedCode int
	}{
		{
			name:         "invalid date_from",
			from:         "2006-01-00gfbfh00",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_from",
			from:         "2006-01-02T15:04:05.0f000000",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_from",
			from:         "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_to",
			to:           "2006-01-00gfbfh00",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_to",
			to:           "invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_to and valid date_form",
			to:           "invalid",
			from:         time.Now().UTC().Format(core.DateStandart),
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_from and valid date_to",
			from:         "invalid",
			to:           time.Now().UTC().Format(core.DateStandart),
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "valid date_from",
			from:         time.Now().UTC().Format(core.DateStandart),
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid date_to",
			to:           time.Now().UTC().Format(core.DateStandart),
			expectedCode: http.StatusOK,
		},
		{
			name:         "all valid",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rec := httptest.NewRecorder()

			url := "/api/v1/logs"

			if tc.from != "" {
				url += "?from=" + tc.from
			}

			if tc.to != "" {
				if tc.from == "" {
					url += "?to=" + tc.to
				}

				if tc.from != "" {
					url += "&to=" + tc.to
				}
			}

			req, _ := http.NewRequest(http.MethodDelete, url, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}

func Test_Server_GetLogRecordsSelectionHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю лог записи
	for i := 0; i < 5; i++ {
		assert.NoError(t, server.TestLogRecord(t, s))
	}

	testCases := []struct {
		name         string
		from         string
		to           string
		page         string
		limit        string
		username     string
		service      int
		expectedCode int
	}{
		{
			name:         "invalid date_from",
			page:         "1",
			limit:        "15",
			from:         "2006-01-00gfbfh00",
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_from",
			page:         "1",
			limit:        "15",
			from:         "2006-01-02T15:04:05.0f000000",
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_from",
			page:         "1",
			limit:        "15",
			from:         "invalid",
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_to",
			page:         "1",
			limit:        "15",
			to:           "2006-01-00gfbfh00",
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_to",
			page:         "1",
			limit:        "15",
			to:           "invalid",
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_to and valid date_form",
			page:         "1",
			limit:        "15",
			to:           "invalid",
			from:         time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid date_from and valid date_to",
			page:         "1",
			limit:        "15",
			from:         "invalid",
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "too big a limit",
			page:         "1",
			limit:        "40",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "page < 1",
			page:         "0",
			limit:        "15",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "limit < 1",
			page:         "1",
			limit:        "0",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid service",
			page:         "1",
			limit:        "15",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      35422,
			username:     "I0HuKc",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid usernmae",
			page:         "1",
			limit:        "15",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "<script></script>",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "invalid usernmae",
			page:         "1",
			limit:        "15",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "GET ALL",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "all valid",
			page:         "1",
			limit:        "15",
			from:         time.Now().UTC().Format(core.DateStandart),
			to:           time.Now().UTC().Format(core.DateStandart),
			service:      AppType.LogTypeAdmin,
			username:     "I0HuKc",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/api/v1/logs?page=" + tc.page

			if tc.limit != "" {
				url += "&limit=" + tc.limit
			}

			if tc.from != "" {
				url += "&from=" + tc.from
			}

			if tc.to != "" {
				url += "&to=" + tc.to
			}

			if strconv.Itoa(tc.service) != "" {
				url += "&service=" + strconv.Itoa(tc.service)
			}

			if tc.username != "" {
				url += "&user=" + tc.username
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	// if rec.Code == http.StatusOK {
	// 	t.Run("checking returned array length", func(t *testing.T) {
	// 		body := map[string]interface{}{}
	// 		assert.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	// 		assert.Len(t, body, 5)
	// 	})
	// }
}

func Test_Server_DeleteLogRecordHandler(t *testing.T) {
	s, redis, teardown := server.TestServer(t)
	defer teardown(redis)

	// Регистрирую менеджера в админке
	tokens, err := server.TestManager(t, s)
	assert.NotNil(t, tokens)
	assert.NoError(t, err)

	// Создаю тестовую лог запись
	assert.NoError(t, server.TestLogRecord(t, s))

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
			id:           "2",
			expectedCode: http.StatusNotFound,
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
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/log/"+tc.id, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
			s.Router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
