package sqlstore_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_LoggerRepository_Delete(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.DB)
	defer teardown("logs")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	lr := &models.LogRecord{
		Info:    "some error",
		Service: AppType.LogLevelServer,
		Module:  "DATABASE",
	}

	assert.NoError(t, s.AdminPanel().Logs().Create(lr))

	testCases := []struct {
		name          string
		id            int
		expectedError bool
	}{
		{
			name:          "valid",
			id:            lr.ID,
			expectedError: false,
		},
		{
			name:          "undefined id",
			id:            10,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lr.ID = tc.id

			if tc.expectedError {
				assert.Error(t, s.AdminPanel().Logs().Delete(lr))
				return
			}

			assert.NoError(t, s.AdminPanel().Logs().Delete(lr))
			assert.NotNil(t, lr)
		})
	}
}

func Test_SQL_LoggerRepository_DeleteSelection(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.DB)
	defer teardown("logs")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	assert.NoError(t, LoggerRepositoryTestCreator(t, s))

	testCases := []struct {
		name,
		from,
		to string
		recreate          func() error
		expectedArrLength int
	}{
		{
			name: "empty all",
			from: "",
			to:   "",
			recreate: func() error {
				return LoggerRepositoryTestCreator(t, s)
			},
			expectedArrLength: 3,
		},
		{
			name: "date_from from future",
			from: time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000"),
			to:   "",
			recreate: func() error {
				return nil
			},
			expectedArrLength: 0,
		},
		{
			name: "date_to from pass",
			from: "",
			to:   "2020-01-02T00:00:00.00000000",
			recreate: func() error {
				return nil
			},
			expectedArrLength: 0,
		},
		{
			name: "all filters are empty",
			from: "",
			to:   "",
			recreate: func() error {
				return nil
			},
			expectedArrLength: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			arr, err := s.AdminPanel().Logs().DeleteSelection(&models.LogRecordSelection{
				DateFrom: tc.from,
				DateTo:   tc.to,
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedArrLength, len(arr))

			assert.NoError(t, tc.recreate())
		})
	}

}

func Test_SQL_LoggerRepository_Selection(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.DB)
	defer teardown("logs")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	assert.NoError(t, LoggerRepositoryTestCreator(t, s))

	u := mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string)

	testCases := []struct {
		name,
		username,
		from,
		service,
		to string
		expectedArrLength int
	}{
		{
			name:              "set username and empty all date and service",
			service:           "",
			from:              "",
			to:                "",
			expectedArrLength: 3,
		},
		{
			name:              "date_from from future",
			username:          u,
			service:           "",
			from:              time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000"),
			to:                "",
			expectedArrLength: 0,
		},
		{
			name:              "date_to from pass",
			username:          u,
			service:           "",
			from:              "",
			to:                "2020-01-02T00:00:00.00000000",
			expectedArrLength: 0,
		},
		{
			name:              "all filters are empty",
			username:          "",
			service:           "",
			from:              "",
			to:                "",
			expectedArrLength: 3,
		},
		{
			name:              "selection by service code",
			username:          "",
			service:           strconv.Itoa(AppType.LogLevelServer),
			from:              "",
			to:                "",
			expectedArrLength: 1,
		},
		{
			name:              "selection by username",
			username:          u,
			service:           "",
			from:              "",
			to:                "",
			expectedArrLength: 1,
		},
		{
			name:              "selection by username and service code",
			username:          u,
			service:           strconv.Itoa(AppType.LogLevelAdmin),
			from:              "",
			to:                "",
			expectedArrLength: 1,
		},
		{
			name:              "selection by username and service code",
			username:          u,
			service:           strconv.Itoa(AppType.LogLevelServer),
			from:              "",
			to:                "",
			expectedArrLength: 0,
		},
	}

	for _, tc := range testCases {
		p, l := 1, 15

		q := &models.LogRecordSelection{
			Page:     &p,
			Limit:    &l,
			Username: tc.username,
			Service:  []string{tc.service},
			DateFrom: tc.from,
			DateTo:   tc.to,
		}

		t.Run("count", func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				c, err := s.AdminPanel().Logs().Count(q)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArrLength, c)
			})
		})

		t.Run("selection", func(t *testing.T) {
			t.Run(tc.name, func(t *testing.T) {
				arr, err := s.AdminPanel().Logs().Selection(q)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArrLength, len(arr))
			})
		})
	}
}

func LoggerRepositoryTestCreator(t *testing.T, s db.SQLStoreI) error {
	t.Helper()

	u := mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string)
	lr := &models.LogRecord{
		Info:    "some error",
		Service: AppType.LogLevelServer,
		Module:  "DATABASE",
	}

	if err := s.AdminPanel().Logs().Create(lr); err != nil {
		return err
	}

	lra := &models.LogRecord{
		Username: &u,
		Info:     "some error",
		Service:  AppType.LogLevelAdmin,
		Module:   "DATABASE",
	}

	if err := s.AdminPanel().Logs().Create(lra); err != nil {
		return err
	}

	lrb := &models.LogRecord{
		Info:    "some error",
		Service: AppType.LogLevelBot,
		Module:  "DATABASE",
	}

	if err := s.AdminPanel().Logs().Create(lrb); err != nil {
		return err
	}

	return nil
}
