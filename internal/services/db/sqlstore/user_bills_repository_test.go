package sqlstore_test

import (
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/sqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_SQL_UserBillsRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	db, teardown := db.TestDB(t, &config.DB)
	defer teardown("users", "bills")

	// Вызываю создание хранилища
	s := sqlstore.Init(db)

	// Регистрация человека как пользователя бота
	u, err := s.User().Create(&models.User{
		ChatID:   int64(mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"].(int)),
		Username: mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string),
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)

	// Создание счета
	nBill, err := s.User().Bills().Create(&models.Bill{
		ChatID: int64(mocks.USER_BILL_REQ["chat_id"].(int)),
		Bill:   mocks.USER_BILL_REQ["bill"].(string),
	})
	assert.NoError(t, err)
	assert.NotNil(t, nBill)

	// Получение списка счетов
	bList, err := s.User().Bills().All(u.ChatID)
	assert.NoError(t, err)
	assert.NotNil(t, bList)
	assert.Len(t, bList, 1)

	// Удаление счета
	dBill, err := s.User().Bills().Delete(nBill)
	assert.NoError(t, err)
	assert.NotNil(t, dBill)

}
