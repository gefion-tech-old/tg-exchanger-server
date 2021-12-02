package db

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

/*
	Интерфейс репозитория для работы с таблицей users
*/
type UserRepository interface {
	/*
		Метод создания новой записи пользователя в таблице users
	*/
	Create(req *models.UserFromBotRequest) (*models.User, error)

	RegisterAsManager(req *models.User) (*models.User, error)

	/*
		Метод обновления записи о пользователе
	*/
	// Update(req *models.UserRequest) (*models.User, error)

	/*
		Метод удаления записи из таблицы users
	*/
	// Delete(chatID int64) (*models.User, error)

	/*
		Метод поиска записи о пользователе в
		таблице users по столбцу chat_id
	*/
	// FindByChatID(chatID int64) (*models.User, error)
}
