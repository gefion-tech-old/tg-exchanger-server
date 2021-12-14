package sqlstore

import (
	"database/sql"
	"time"

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

func (r *ExchangerRepository) Create(e *models.Exchanger) (*models.Exchanger, error) {
	if err := r.store.QueryRow(
		`
		INSERT INTO exchangers(name, url)
		SELECT $1, $2
		RETURNING id, name, url, created_at, updated_at
		`,
		e.Name,
		e.UrlToParse,
	).Scan(
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *ExchangerRepository) Update(e *models.Exchanger) (*models.Exchanger, error) {
	if err := r.store.QueryRow(
		`
		UPDATE exchangers
		SET name=$1, url=$2, updated_at=$3
		WHERE id=$4
		RETURNING id, name, url, created_at, updated_at
		`,
		e.Name,
		e.UrlToParse,
		time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
	).Scan(
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *ExchangerRepository) Count() (int, error) {
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

func (r *ExchangerRepository) Get(e *models.Exchanger) (*models.Exchanger, error) {
	if err := r.store.QueryRow(
		`
		SELECT id, name, url, created_at, updated_at
		FROM exchangers
		WHERE id=$1
		`,
		e.ID,
	).Scan(
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *ExchangerRepository) GetSlice(limit int) ([]*models.Exchanger, error) {
	eArr := []*models.Exchanger{}

	rows, err := r.store.Query(
		`
		SELECT id, name, url, created_at, updated_at
		FROM exchangers
		ORDER BY id DESC
		LIMIT $1
		`,
		limit,
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
			&e.UpdatedAt,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			continue
		}

		eArr = append(eArr, &e)

	}

	return eArr, nil

}

func (r *ExchangerRepository) Delete(e *models.Exchanger) (*models.Exchanger, error) {
	if err := r.store.QueryRow(
		`
		DELETE FROM exchangers
		WHERE id=$1
		RETURNING id, name, url, created_at, updated_at
		`,
		e.ID,
	).Scan(
		&e.ID,
		&e.Name,
		&e.UpdatedAt,
		&e.CreatedAt,
		&e.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return e, nil
}
