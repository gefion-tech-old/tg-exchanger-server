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

func Test_SQL_UserRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	db, teardown := db.TestDB(t, &config.Services.DB)
	defer teardown("users")

	// Вызываю создание хранилища
	s := sqlstore.Init(db)

	u := &models.User{
		ChatID:   int64(mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"].(int)),
		Username: mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string),
	}

	// Регистрация человека как пользователя бота
	err := s.User().Create(u)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	// Регистрация человека как менеджера
	assert.NoError(t, s.User().RegisterInAdminPanel(u))

	// Поиск пользователя по его username
	uUsername, err := s.User().FindByUsername(u.Username)
	assert.NoError(t, err)
	assert.NotNil(t, uUsername)
}
