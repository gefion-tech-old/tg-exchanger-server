package sqlstore_test

import (
	"testing"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_LoggerRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.DB)
	defer teardown("logs")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	u := mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string)
	lr := &models.LogRecord{
		Username: &u,
		Info:     "some error",
		Service:  static.L__SERVER,
		Module:   "DATABASE",
	}

	assert.NoError(t, s.AdminPanel().Logs().Create(lr))
	assert.NotNil(t, lr)

	testCases := []struct {
		name,
		username,
		from,
		to string
		expectedArrLength int
	}{
		{
			name:              "fill username and empty all date",
			username:          u,
			from:              "",
			to:                "",
			expectedArrLength: 1,
		},

		{
			name:              "date_from from future",
			username:          u,
			from:              time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000"),
			to:                "",
			expectedArrLength: 0,
		},

		{
			name:              "date_to from pass",
			username:          u,
			from:              "",
			to:                "2020-01-02T00:00:00.00000000",
			expectedArrLength: 0,
		},

		{
			name:              "all filters are empty",
			username:          "",
			from:              "",
			to:                "",
			expectedArrLength: 1,
		},
	}

	for _, tc := range testCases {
		t.Run("count", func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				c, err := s.AdminPanel().Logs().CountWithCustomFilters(tc.username, tc.from, tc.to)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArrLength, c)
			})
		})

		t.Run("selection", func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				arr, err := s.AdminPanel().Logs().SelectionWithCustomFilters(1, 15, tc.username, tc.from, tc.to)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArrLength, len(arr))
			})
		})

		t.Run("delete selection", func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				arr, err := s.AdminPanel().Logs().DeleteSelection(tc.from, tc.to)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArrLength, len(arr))

				if len(arr) > 0 {
					t.Run("create record after delete selection", func(t *testing.T) {
						assert.NoError(t, s.AdminPanel().Logs().Create(lr))
						assert.NotNil(t, lr)
					})
				}
			})
		})
	}

	assert.NoError(t, s.AdminPanel().Logs().Delete(lr))
	assert.NotNil(t, lr)
}
