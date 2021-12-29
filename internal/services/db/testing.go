package db

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/mocks"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

func TestDB(t *testing.T, config *config.DatabaseConfig) (*sql.DB, func(...string)) {
	t.Helper()

	// Инициализирую БД
	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		t.Fatal(err)
	}

	// Поверка подключения к БД
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		// Проверяю, переданы ли таблицы которые необходимо очистить
		if len(tables) > 0 {
			// Очищаю таблицы
			str := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")) // Создаю строку команды для очистки бд
			db.Exec(str)
		}

		db.Close()
	}
}

func CreateUser(t *testing.T, s SQLStoreI) (*models.User, error) {
	t.Helper()

	u := &models.User{
		ChatID:   int64(mocks.USER_IN_BOT_REGISTRATION_REQ["chat_id"].(int)),
		Username: mocks.USER_IN_BOT_REGISTRATION_REQ["username"].(string),
	}

	// Регистрация человека как пользователя бота
	err := s.User().Create(u)
	if err != nil {
		return nil, err
	}

	// Регистрация человека как менеджера
	if err := s.User().RegisterInAdminPanel(u); err != nil {
		return nil, err
	}

	return u, nil
}
