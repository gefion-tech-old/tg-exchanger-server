package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserBillsRepository struct {
	store *sql.DB
}

/*
	Добавить банковский счет пользователю
*/
func (r *UserBillsRepository) Create(b *models.Bill) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO bills(chat_id, bill)
		SELECT $1, $2
		WHERE NOT EXISTS (SELECT chat_id FROM bills WHERE chat_id=$3 AND bill=$4)
		RETURNING id, chat_id, bill, created_at
		`,
		b.ChatID,
		b.Bill,
		b.ChatID,
		b.Bill,
	).Scan(
		&b.ID,
		&b.ChatID,
		&b.Bill,
		&b.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

/*
	Удалить банковский счет пользователя
*/
func (r *UserBillsRepository) Delete(b *models.Bill) error {
	if err := r.store.QueryRow(
		`
		DELETE FROM bills
		WHERE chat_id=$1 AND id=$2
		RETURNING id, chat_id, bill, created_at
		`,
		b.ChatID,
		b.ID,
	).Scan(
		&b.ID,
		&b.ChatID,
		&b.Bill,
		&b.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (r *UserBillsRepository) FindById(b *models.Bill) error {
	if err := r.store.QueryRow(
		`
		SELECT id, chat_id, bill, created_at 
		FROM bills WHERE id=$1
		`,
		b.ID,
	).Scan(
		&b.ID,
		&b.ChatID,
		&b.Bill,
		&b.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

/*
	Получить все счета пользователя
*/
func (r *UserBillsRepository) All(chatID int64) ([]*models.Bill, error) {
	ub := []*models.Bill{}

	rows, err := r.store.Query(
		`
		SELECT id, chat_id, bill, created_at 
		FROM bills WHERE chat_id=$1
		`,
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		b := &models.Bill{}
		if err := rows.Scan(
			&b.ID,
			&b.ChatID,
			&b.Bill,
			&b.CreatedAt,
		); err != nil {
			continue
		}

		ub = append(ub, b)
	}

	return ub, nil
}
