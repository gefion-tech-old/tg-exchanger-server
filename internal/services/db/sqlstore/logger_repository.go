package sqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
)

type LoggerRepository struct {
	store *sql.DB
}

// # TESTED
func (r *LoggerRepository) Create(lr *models.LogRecord) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO logs (username, info, service, module)
		SELECT $1, $2, $3, $4	
		RETURNING id, username, info, service, module, created_at
		`,
		lr.Username,
		lr.Info,
		lr.Service,
		lr.Module,
	).Scan(
		&lr.ID,
		&lr.Username,
		&lr.Info,
		&lr.Service,
		&lr.Module,
		&lr.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

// # TESTED
func (r *LoggerRepository) Delete(lr *models.LogRecord) error {
	if err := r.store.QueryRow(
		`
		DELETE FROM logs
		WHERE id=$1
		RETURNING id, username, info, service, module, created_at
		`,
		lr.ID,
	).Scan(
		&lr.ID,
		&lr.Username,
		&lr.Info,
		&lr.Service,
		&lr.Module,
		&lr.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

// # TESTED
func (r *LoggerRepository) CountWithCustomFilters(username, date_from, date_to string) (int, error) {
	var c int

	if date_from == "" {
		date_from = "2020-01-01T00:00:00.00000000"
	}

	if date_to == "" {
		date_to = time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000")

	}

	if username != "" {
		if err := r.store.QueryRow(
			`
			SELECT count(*)
			FROM logs
			WHERE created_at >= $1 AND created_at < $2 AND username=$3	
			`,
			date_from,
			date_to,
			username,
		).Scan(
			&c,
		); err != nil {
			return 0, err
		}
	}

	if err := r.store.QueryRow(
		`
		SELECT count(*)
		FROM logs
		WHERE created_at >= $1 AND created_at < $2 			
		`,
		date_from,
		date_to,
	).Scan(
		&c,
	); err != nil {
		return 0, err
	}

	return c, nil
}

// # TESTED
func (r *LoggerRepository) SelectionWithCustomFilters(page, limit int, username, date_from, date_to string) ([]*models.LogRecord, error) {
	arr := []*models.LogRecord{}
	var rows *sql.Rows

	if date_from == "" {
		date_from = "2020-01-01T00:00:00.00000000"
	}

	if date_to == "" {
		date_to = time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000")

	}

	if username == "" {
		r, err := r.store.Query(
			`
			SELECT id, username, info, service, module, created_at
			FROM logs
			WHERE created_at >= $1 AND created_at < $2
			ORDER BY id DESC
			OFFSET $3
			LIMIT $4
			`,
			date_from,
			date_to,
			tools.OffsetThreshold(page, limit),
			limit,
		)
		if err != nil {
			return nil, err
		}

		rows = r
		defer r.Close()
	}

	if username != "" {
		r, err := r.store.Query(
			`
			SELECT id, username, info, service, module, created_at
			FROM logs
			WHERE created_at >= $1 AND created_at < $2 AND username=$3
			ORDER BY id DESC
			OFFSET $4
			LIMIT $5
			`,
			date_from,
			date_to,
			username,
			tools.OffsetThreshold(page, limit),
			limit,
		)
		if err != nil {
			return nil, err
		}

		rows = r
		defer r.Close()
	}

	if rows != nil {
		for rows.Next() {
			lr := &models.LogRecord{}
			if err := rows.Scan(
				&lr.ID,
				&lr.Username,
				&lr.Info,
				&lr.Service,
				&lr.Module,
				&lr.CreatedAt,
			); err != nil {
				continue
			}

			arr = append(arr, lr)
		}

		return arr, nil
	}

	return nil, errors.ErrInvalidCondition
}

// # TESTED
func (r *LoggerRepository) DeleteSelection(date_from, date_to string) ([]*models.LogRecord, error) {
	arr := []*models.LogRecord{}

	if date_from == "" {
		date_from = "2020-01-01T00:00:00.00000000"
	}

	if date_to == "" {
		date_to = time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000")

	}

	rows, err := r.store.Query(
		`
		DELETE FROM logs
		WHERE created_at >= $1 AND created_at < $2	
		RETURNING id, username, info, service, module, created_at
		`,
		date_from,
		date_to,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		lr := &models.LogRecord{}
		if err := rows.Scan(
			&lr.ID,
			&lr.Username,
			&lr.Info,
			&lr.Service,
			&lr.Module,
			&lr.CreatedAt,
		); err != nil {
			continue
		}

		arr = append(arr, lr)
	}

	return arr, nil
}
