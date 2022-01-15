package sqlstore_test

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_NotificationRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	db, teardown := db.TestDB(t, &config.Services.DB)
	defer teardown("notifications")

	// Вызываю создание хранилища
	s := sqlstore.Init(db)

	var data_854 *models.Notification
	mapstructure.Decode(mocks.ADMIN_NOTIFICATION_854, &data_854)

	// Создание уведомления 854
	assert.NoError(t, s.AdminPanel().Notification().Create(data_854))

	// Поиск созданного уведомления
	assert.NoError(t, s.AdminPanel().Notification().Get(data_854))

	q := &models.NotificationSelection{
		Page:  1,
		Limit: 10,
	}

	// Получить уведомления из БД
	arrN, err := s.AdminPanel().Notification().Selection(q)
	assert.NoError(t, err)
	assert.NotNil(t, arrN)
	assert.Len(t, arrN, 1)

	// Подсчет кол-ва
	c, err := s.AdminPanel().Notification().Count(nil)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, len(arrN), c)

	// Удалить запись
	assert.NoError(t, s.AdminPanel().Notification().Delete(data_854))

	// Получить запись
	assert.Error(t, s.AdminPanel().Notification().Get(data_854))

}
