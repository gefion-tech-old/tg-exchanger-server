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

func Test_SQL_BotMessagesRepository(t *testing.T) {
	config := config.InitTestConfig(t)

	database, teardown := db.TestDB(t, &config.DB)
	defer teardown("users", "bot_messages")

	// Вызываю создание хранилища
	s := sqlstore.Init(database)

	u, err := db.CreateUser(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	var m *models.BotMessage
	mapstructure.Decode(mocks.BOT_MESSAGE_REQ, &m)
	m.MessageText = "some text"
	m.CreatedBy = mocks.MANAGER_IN_ADMIN_REQ["username"].(string)

	// Создание сообщения
	assert.NoError(t, s.AdminPanel().BotMessages().Create(m))
	assert.NotNil(t, m)

	// Получить из БД
	assert.NoError(t, s.AdminPanel().BotMessages().Get(m))
	assert.NotNil(t, m)

	// Обновление
	m.MessageText = "new text"
	assert.NoError(t, s.AdminPanel().BotMessages().Update(m))
	assert.NotNil(t, m)

	// Удаление
	assert.NoError(t, s.AdminPanel().BotMessages().Delete(m))
	assert.NotNil(t, m)
}
