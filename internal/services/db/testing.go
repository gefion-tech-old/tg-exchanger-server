package db

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
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
