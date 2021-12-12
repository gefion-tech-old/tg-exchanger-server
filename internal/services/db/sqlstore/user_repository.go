package sqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type UserRepository struct {
	store               *sql.DB
	userBillsRepository *UserBillsRepository
}

/*
	==========================================================================================
	КОНСТРУКТОРЫ ВЛОЖЕННЫХ СТРУКТУР
	==========================================================================================
*/

func (r *UserRepository) Bills() db.UserBillsRepository {
	if r.userBillsRepository != nil {
		return r.userBillsRepository
	}

	r.userBillsRepository = &UserBillsRepository{
		store: r.store,
	}

	return r.userBillsRepository
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ ТЕКУЩЕЙ СТРУКТУРЫ
	==========================================================================================
*/

/*
	Метод создания новой записи пользователя в таблице users
*/
func (r *UserRepository) Create(u *models.User) (*models.User, error) {
	if err := r.store.QueryRow(
		`
		INSERT INTO users (chat_id, username) 
		SELECT $1, $2
		WHERE NOT EXISTS (SELECT chat_id FROM users WHERE chat_id=$3)
		RETURNING chat_id, username, hash, created_at, updated_at
		`,
		u.ChatID,
		u.Username,
		u.ChatID,
	).Scan(
		&u.ChatID,
		&u.Username,
		&u.Hash,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}

/*
	Метод регистрации человека как менеджера для доступа к админке
*/
func (r *UserRepository) RegisterAsManager(u *models.User) (*models.User, error) {
	if err := r.store.QueryRow(
		`
		UPDATE users
		SET hash=$1, updated_at=$2
		WHERE username=$3
		RETURNING chat_id, username, hash, created_at, updated_at
		`,
		u.Hash,
		time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
		u.Username,
	).Scan(
		&u.ChatID,
		&u.Username,
		&u.Hash,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}

/*
	Метод поиска записи о пользователе в таблице users по столбцу username
*/
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	u := &models.User{}

	if err := r.store.QueryRow(
		`
		SELECT chat_id, hash, username, created_at, updated_at 
		FROM users 
		WHERE username=$1
		`,
		username,
	).Scan(
		&u.ChatID,
		&u.Hash,
		&u.Username,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetAllManagers() ([]*models.User, error) {
	uArr := []*models.User{}

	rows, err := r.store.Query(
		`
		SELECT chat_id, hash, username, created_at, updated_at
		FROM users WHERE hash IS NOT NULL
		`,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := models.User{}
		if err := rows.Scan(
			&u.ChatID,
			&u.Hash,
			&u.Username,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			continue
		}

		uArr = append(uArr, &u)
	}
	defer rows.Close()

	return uArr, nil
}
