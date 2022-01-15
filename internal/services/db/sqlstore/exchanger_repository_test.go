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

func Test_SQL_ExcangerRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	db, teardown := db.TestDB(t, &config.Services.DB)
	defer teardown("exchangers")

	// Вызываю создание хранилища
	s := sqlstore.Init(db)

	var data *models.Exchanger
	mapstructure.Decode(mocks.ADMIN_EXCHANGER, &data)

	// Создание
	assert.NoError(t, s.AdminPanel().Exchanger().Create(data))

	// Обновление
	data.Name = "new"
	assert.NoError(t, s.AdminPanel().Exchanger().Update(data))

	q := &models.ExchangerSelection{
		Page:  1,
		Limit: 10,
	}

	// Подсчет
	c, err := s.AdminPanel().Exchanger().Count(q)
	assert.NoError(t, err)
	assert.Equal(t, 1, c)

	// Получение среза
	slice, err := s.AdminPanel().Exchanger().Selection(q)
	assert.NoError(t, err)
	assert.Len(t, slice, 1)

	// Удаление
	assert.NoError(t, s.AdminPanel().Exchanger().Delete(data))
}
