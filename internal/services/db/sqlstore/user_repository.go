package sqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserRepository struct {
	store *sql.DB
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ ТЕКУЩЕЙ СТРУКТУРЫ
	==========================================================================================
*/

/*
	Метод создания новой записи пользователя в таблице users
*/
func (r *UserRepository) Create(u *models.User) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO users (chat_id, username) 
		SELECT $1, $2
		WHERE NOT EXISTS (SELECT chat_id FROM users WHERE chat_id=$3 OR username=$4)
		RETURNING chat_id, username, hash, role, created_at, updated_at
		`,
		u.ChatID,
		u.Username,
		u.ChatID,
		u.Username,
	).Scan(
		&u.ChatID,
		&u.Username,
		&u.Hash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

/*
	Метод регистрации человека как менеджера для доступа к админке
*/
func (r *UserRepository) RegisterInAdminPanel(u *models.User) error {
	if err := r.store.QueryRow(
		`
		UPDATE users
		SET hash=$1, role=$2, updated_at=$3
		WHERE username=$4
		RETURNING chat_id, username, hash, role, created_at, updated_at
		`,
		u.Hash,
		u.Role,
		time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
		u.Username,
	).Scan(
		&u.ChatID,
		&u.Username,
		&u.Hash,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

/*
	Метод поиска записи о пользователе в таблице users по столбцу username
*/
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	u := &models.User{}

	if err := r.store.QueryRow(
		`
		SELECT chat_id, username, hash, role, created_at, updated_at 
		FROM users 
		WHERE username=$1
		`,
		username,
	).Scan(
		&u.ChatID,
		&u.Username,
		&u.Hash,
		&u.Role,
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
		SELECT chat_id, username, hash, role, created_at, updated_at
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
			&u.Username,
			&u.Hash,
			&u.Role,
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
