package sqlstore_test

import (
	"fmt"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_MerchantAutopayoutRepository_Selection_Count(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.Services.DB)
	defer teardown("users", "bot_messages", "merchant_autopayout")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	u, err := db.CreateUser(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	// Создание сообщения
	var m *models.BotMessage
	assert.NoError(t, mapstructure.Decode(mocks.BOT_MESSAGE_REQ, &m))
	m.MessageText = "some text"
	m.CreatedBy = mocks.MANAGER_IN_ADMIN_REQ["username"].(string)
	assert.NoError(t, s.AdminPanel().BotMessages().Create(m))
	assert.NotNil(t, m)

	var ma *models.MerchantAutopayout
	assert.NoError(t, mapstructure.Decode(mocks.MerchantAutopayout, &ma))
	ma.MessageID = m.ID
	ma.CreatedBy = u.Username

	for i := 0; i < 6; i++ {
		ma.Name = fmt.Sprintf("%d#m", i)
		assert.NoError(t, s.AdminPanel().MerchantAutopayout().Create(ma))
	}

	for i := 0; i < 10; i++ {
		ma.Name = fmt.Sprintf("%d#a", i)
		ma.Service = AppType.MerchantAutoPayoutMine
		assert.NoError(t, s.AdminPanel().MerchantAutopayout().Create(ma))
	}

	t.Run("Count", func(t *testing.T) {
		testCases := []struct {
			name              string
			paload            func() interface{}
			expectedArrLength int
		}{
			{
				name: "1",
				paload: func() interface{} {
					p, l := 1, 5
					return &models.MerchantAutopayoutSelection{
						Page:    &p,
						Limit:   &l,
						Service: []string{AppType.MerchantAutoPayoutWhitebit},
					}
				},
				expectedArrLength: 5,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, err := s.AdminPanel().MerchantAutopayout().Count(tc.paload())
				assert.NoError(t, err)
				assert.Equal(t, 6, c)
			})
		}

	})

	t.Run("Selection", func(t *testing.T) {
		testCases := []struct {
			name              string
			paload            func() interface{}
			expectedArrLength int
		}{
			{
				name: "1",
				paload: func() interface{} {
					p, l := 1, 5
					return &models.MerchantAutopayoutSelection{
						Page:    &p,
						Limit:   &l,
						Service: []string{AppType.MerchantAutoPayoutWhitebit},
					}
				},
				expectedArrLength: 5,
			},
			{
				name: "2",
				paload: func() interface{} {
					p, l := 1, 15
					return &models.MerchantAutopayoutSelection{
						Page:    &p,
						Limit:   &l,
						Service: []string{AppType.MerchantAutoPayoutWhitebit},
					}
				},
				expectedArrLength: 6,
			},
			{
				name: "3",
				paload: func() interface{} {
					p, l := 1, 1
					return &models.MerchantAutopayoutSelection{
						Page:    &p,
						Limit:   &l,
						Service: []string{AppType.MerchantAutoPayoutWhitebit},
					}
				},
				expectedArrLength: 1,
			},
			{
				name: "4",
				paload: func() interface{} {
					p, l := 1, 15
					return &models.MerchantAutopayoutSelection{
						Page:    &p,
						Limit:   &l,
						Service: []string{AppType.MerchantAutoPayoutMine},
					}
				},
				expectedArrLength: 10,
			},
			{
				name: "5",
				paload: func() interface{} {
					p, l := 2, 9
					return &models.MerchantAutopayoutSelection{
						Page:    &p,
						Limit:   &l,
						Service: []string{AppType.MerchantAutoPayoutMine},
					}
				},
				expectedArrLength: 1,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				arr, err := s.AdminPanel().MerchantAutopayout().Selection(tc.paload())
				assert.NoError(t, err)
				assert.Len(t, arr, tc.expectedArrLength)
			})
		}
	})
}

func Test_SQL_MerchantAutopayoutRepository_Create_Update_Delete(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.Services.DB)
	defer teardown("users", "bot_messages", "merchant_autopayout")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	u, err := db.CreateUser(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	// Создание сообщения
	var m *models.BotMessage
	assert.NoError(t, mapstructure.Decode(mocks.BOT_MESSAGE_REQ, &m))
	m.MessageText = "some text"
	m.CreatedBy = mocks.MANAGER_IN_ADMIN_REQ["username"].(string)
	assert.NoError(t, s.AdminPanel().BotMessages().Create(m))
	assert.NotNil(t, m)

	var ma *models.MerchantAutopayout
	assert.NoError(t, mapstructure.Decode(mocks.MerchantAutopayout, &ma))
	ma.MessageID = m.ID
	ma.CreatedBy = u.Username

	t.Run("Create", func(t *testing.T) {
		assert.NoError(t, s.AdminPanel().MerchantAutopayout().Create(ma))
		assert.NotNil(t, ma)
	})

	t.Run("Update", func(t *testing.T) {
		testCases := []struct {
			name    string
			payload func() *models.MerchantAutopayout
		}{
			{
				name: "update name",
				payload: func() *models.MerchantAutopayout {
					newMa := ma
					newMa.Name = "new_name"
					return newMa
				},
			},
			{
				name: "update service_type",
				payload: func() *models.MerchantAutopayout {
					newMa := ma
					newMa.ServiceType = AppType.UseAsAutoPayout
					return newMa
				},
			},
			{
				name: "update options",
				payload: func() *models.MerchantAutopayout {
					newMa := ma
					newMa.Options = `{"test": 1}`
					return newMa
				},
			},
			{
				name: "update options",
				payload: func() *models.MerchantAutopayout {
					newMa := ma
					newMa.Options = `{"test": 1}`
					return newMa
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				p := tc.payload()
				assert.NoError(t, s.AdminPanel().MerchantAutopayout().Update(p))
				assert.NotNil(t, ma)
				assert.Equal(t, p, ma)
			})
		}
	})

	t.Run("Delete", func(t *testing.T) {
		assert.NoError(t, s.AdminPanel().MerchantAutopayout().Delete(ma))
		assert.NotNil(t, ma)
	})
}
