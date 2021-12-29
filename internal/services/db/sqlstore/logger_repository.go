package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
)

type LoggerRepository struct {
	store *sql.DB
}

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

func (r *LoggerRepository) Count() (int, error) {
	var c int
	if err := r.store.QueryRow(
		`
		SELECT count(*)
		FROM logs		
		`,
	).Scan(
		&c,
	); err != nil {
		return 0, err
	}

	return c, nil
}

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

func (r *LoggerRepository) Selection(page, limit int) ([]*models.LogRecord, error) {
	arr := []*models.LogRecord{}

	rows, err := r.store.Query(
		`
		SELECT id, username, info, service, module, created_at
		FROM logs
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
