package db

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

/*
	Интерфейс репозитория для работы с таблицей users
*/
type UserRepository interface {
	Bills() UserBillsRepository
	/*
		Метод создания новой записи пользователя в таблице users
	*/
	Create(req *models.UserFromBotRequest) (*models.User, error)

	/*
		Метод регистрации человека как менеджера для доступа к админке
	*/
	RegisterAsManager(req *models.User) (*models.User, error)

	/*
		Метод поиска записи о пользователе в
		таблице users по столбцу username
	*/
	FindByUsername(username string) (*models.User, error)
}

type UserBillsRepository interface {
	/*
		Добавить банковский счет пользователю
	*/
	Create(b *models.Bill) (*models.Bill, error)

	/*
		Удалить банковский счет пользователя
	*/
	Delete(b *models.Bill) (*models.Bill, error)

	/*
		Получить все счета пользователя
	*/
	All(chatID int64) ([]*models.Bill, error)
}
