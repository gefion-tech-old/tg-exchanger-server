package sqlstore_test

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_UserBillsRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	db, teardown := db.TestDB(t, &config.Services.DB)
	defer teardown("users", "bills")

	// Вызываю создание хранилища
	s := sqlstore.Init(db)

	u := models.User{
		ChatID:   int64(mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"].(int)),
		Username: mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string),
	}

	// Регистрация человека как пользователя бота
	err := s.User().Create(&u)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	b := &models.Bill{
		ChatID: int64(mocks.USER_BILL_REQ["chat_id"].(int)),
		Bill:   mocks.USER_BILL_REQ["bill"].(string),
	}

	// Создание счета
	assert.NoError(t, s.AdminPanel().Bills().Create(b))
	assert.NotNil(t, b)

	// Получение одного счета
	assert.NoError(t, s.AdminPanel().Bills().FindById(&models.Bill{ID: b.ID}))
	assert.NotNil(t, b)

	// Получение списка счетов
	bList, err := s.AdminPanel().Bills().All(u.ChatID)
	assert.NoError(t, err)
	assert.NotNil(t, bList)
	assert.Len(t, bList, 1)

	// Удаление счета
	assert.NoError(t, s.AdminPanel().Bills().Delete(b))
}
