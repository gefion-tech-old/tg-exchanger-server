package sqlstore_test

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_NotificationRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	db, teardown := db.TestDB(t, &config.DB)
	defer teardown("notifications")

	// Вызываю создание хранилища
	s := sqlstore.Init(db)

	var data_854 *models.Notification
	mapstructure.Decode(mocks.ADMIN_NOTIFICATION_854, &data_854)

	// Создание уведомления 854
	n, err := s.Manager().Notification().Create(data_854)
	assert.NoError(t, err)
	assert.NotNil(t, n)

	// Поиск созданного уведомления
	n2, err := s.Manager().Notification().Get(n)
	assert.NoError(t, err)
	assert.NotNil(t, n2)

	// Получить уведомления из БД
	arrN, err := s.Manager().Notification().Selection(1, 10)
	assert.NoError(t, err)
	assert.NotNil(t, arrN)
	assert.Len(t, arrN, 1)

	// Подсчет кол-ва
	c, err := s.Manager().Notification().Count()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, len(arrN), c)

	// Удалить запись
	d, err := s.Manager().Notification().Delete(n)
	assert.NoError(t, err)
	assert.NotNil(t, d)
	_, dErr := s.Manager().Notification().Get(n)
	assert.Error(t, dErr)

}
