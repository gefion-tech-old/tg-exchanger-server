package sqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
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
func (r *ExchangerRepository) Create(e *models.Exchanger) (*models.Exchanger, error) {
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
		return nil, err
	}

	return e, nil
}

/*
	Обновить запись в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) Update(e *models.Exchanger) (*models.Exchanger, error) {
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
		return nil, err
	}

	return e, nil
}

/*
	Подсчет кол-ва записей в таблице `exchangers`

	# TESTED
*/
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

func (r *ExchangerRepository) Selection(page, limit int) ([]*models.Exchanger, error) {
	arr := []*models.Exchanger{}

	rows, err := r.store.Query(
		`
		SELECT id, name, url, created_by, created_at, updated_at
		FROM exchangers
		ORDER BY id DESC
		OFFSET $1
		LIMIT $2
		`,
		tools.OffsetThreshold(page, limit),
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
func (r *ExchangerRepository) GetByName(e *models.Exchanger) (*models.Exchanger, error) {
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
		return nil, err
	}

	return e, nil
}

/*
	Получить лдимитированный объем записей в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) GetSlice(limit int) ([]*models.Exchanger, error) {
	eArr := []*models.Exchanger{}

	rows, err := r.store.Query(
		`
		SELECT id, name, url, created_by, created_at, updated_at
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
			&e.UrlToParse,
			&e.CreatedBy,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			continue
		}

		eArr = append(eArr, &e)

	}

	return eArr, nil
}

/*
	Удалить запись в таблице `exchangers`

	# TESTED
*/
func (r *ExchangerRepository) Delete(e *models.Exchanger) (*models.Exchanger, error) {
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
		return nil, err
	}

	return e, nil
}
