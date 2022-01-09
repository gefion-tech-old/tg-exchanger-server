package sqlstore

import (
	"database/sql"
	"time"

	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type ExchangerRepository struct {
	store *sql.DB
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ ТЕКУЩЕЙ СТРУКТУРЫ
	==========================================================================================
*/

/*
	Создать запись в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) Create(e *models.Exchanger) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO exchangers(name, url, created_by)
		SELECT $1, $2, $3
		RETURNING id, name, url, created_by, created_at, updated_at
		`,
		e.Name,
		e.UrlToParse,
		e.CreatedBy,
	).Scan(
		&e.ID,
		&e.Name,
		&e.UrlToParse,
		&e.CreatedBy,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

/*
	Обновить запись в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) Update(e *models.Exchanger) error {
	if err := r.store.QueryRow(
		`
		UPDATE exchangers
		SET name=$1, url=$2, updated_at=$3
		WHERE id=$4
		RETURNING id, name, url, created_by, created_at, updated_at
		`,
		e.Name,
		e.UrlToParse,
		time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
		e.ID,
	).Scan(
		&e.ID,
		&e.Name,
		&e.UrlToParse,
		&e.CreatedBy,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

/*
	Подсчет кол-ва записей в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) Count(querys interface{}) (int, error) {
	var c int
	if err := r.store.QueryRow(
		`
		SELECT count(*)
		FROM exchangers		
		`,
	).Scan(
		&c,
	); err != nil {
		return 0, err
	}

	return c, nil
}

func (r *ExchangerRepository) Selection(querys interface{}) ([]*models.Exchanger, error) {
	arr := []*models.Exchanger{}
	q := querys.(*models.ExchangerSelection)

	rows, err := r.store.Query(
		`
		SELECT id, name, url, created_by, created_at, updated_at
		FROM exchangers
		ORDER BY id DESC
		OFFSET $1
		LIMIT $2
		`,
		AppMath.OffsetThreshold(q.Page, q.Limit),
		q.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := models.Exchanger{}
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.UrlToParse,
			&e.CreatedBy,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			continue
		}

		arr = append(arr, &e)

	}

	return arr, nil
}

/*
	Поиск записи по поля `name` в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) GetByName(e *models.Exchanger) error {
	if err := r.store.QueryRow(
		`
		SELECT id, name, url, created_by, created_at, updated_at
		FROM exchangers
		WHERE name=$1		
		`,
		e.Name,
	).Scan(
		&e.ID,
		&e.Name,
		&e.UrlToParse,
		&e.CreatedBy,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

/*
	Удалить запись в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) Delete(e *models.Exchanger) error {
	if err := r.store.QueryRow(
		`
		DELETE FROM exchangers
		WHERE id=$1
		RETURNING id, name, url, created_by, created_at, updated_at
		`,
		e.ID,
	).Scan(
		&e.ID,
		&e.Name,
		&e.UrlToParse,
		&e.CreatedBy,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}
