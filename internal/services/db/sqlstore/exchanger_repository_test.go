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

func Test_SQL_ExcangerRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	db, teardown := db.TestDB(t, &config.DB)
	defer teardown("exchangers")

	// Вызываю создание хранилища
	s := sqlstore.Init(db)

	var data *models.Exchanger
	mapstructure.Decode(mocks.ADMIN_EXCHANGER, &data)

	// Создание
	e1, err := s.Manager().Exchanger().Create(data)
	assert.NoError(t, err)
	assert.NotNil(t, e1)

	// Обновление
	data.Name = "new"
	e2, err := s.Manager().Exchanger().Update(data)
	assert.NoError(t, err)
	assert.NotNil(t, e2)

	// Подсчет
	c, err := s.Manager().Exchanger().Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, c)

	// Получение среза
	slice, err := s.Manager().Exchanger().GetSlice(10)
	assert.NoError(t, err)
	assert.Len(t, slice, 1)

	// Удаление
	e3, err := s.Manager().Exchanger().Delete(data)
	assert.NoError(t, err)
	assert.NotNil(t, e3)
}
