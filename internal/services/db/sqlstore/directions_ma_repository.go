package sqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type DirectionsMaRepository struct {
	store *sql.DB
}

func (r *DirectionsMaRepository) Create(dma *models.DirectionMA) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO directions_ma(direction_id, ma_id, service_type, status)
		SELECT $1, $2, $3, $4
		RETURNING id, direction_id, ma_id, service_type, status, created_at, updated_at
		`,
		dma.DirectionID,
		dma.MaID,
		dma.ServiceType,
		dma.Status,
	).Scan(
		&dma.ID,
		&dma.DirectionID,
		&dma.MaID,
		&dma.ServiceType,
		&dma.Status,
		&dma.CreatedAt,
		&dma.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsMaRepository) Update(dma *models.DirectionMA) error {
	if err := r.store.QueryRow(
		`
		UPDATE directions_ma
		SET service_type=$1, status=$2, updated_at=$3
		WHERE id=$4
		RETURNING id, direction_id, ma_id, service_type, status, created_at, updated_at
		`,
		dma.ServiceType,
		dma.Status,
		time.Now().UTC().Format(core.DateStandart),
		dma.ID,
	).Scan(
		&dma.ID,
		&dma.DirectionID,
		&dma.MaID,
		&dma.ServiceType,
		&dma.Status,
		&dma.CreatedAt,
		&dma.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsMaRepository) Get(dma *models.DirectionMA) error {
	if err := r.store.QueryRow(
		`
		SELECT id, direction_id, ma_id, service_type, status, created_at, updated_at
		FROM directions_ma
		WHERE id=$1
		`,
		dma.ID,
	).Scan(
		&dma.ID,
		&dma.DirectionID,
		&dma.MaID,
		&dma.ServiceType,
		&dma.Status,
		&dma.CreatedAt,
		&dma.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsMaRepository) Delete(dma *models.DirectionMA) error {
	if err := r.store.QueryRow(
		`
		DELETE FROM directions_ma
		WHERE id=$1
		RETURNING id, direction_id, ma_id, service_type, status, created_at, updated_at
		`,
		dma.ID,
	).Scan(
		&dma.ID,
		&dma.DirectionID,
		&dma.MaID,
		&dma.ServiceType,
		&dma.Status,
		&dma.CreatedAt,
		&dma.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsMaRepository) Count(querys interface{}) (int, error) {
	q := querys.(*models.DirectionMASelection)
	var c int

	if err := r.store.QueryRow(
		`
		SELECT count(*)
		FROM directions_ma
		WHERE direction_id=$1
		`,
		q.DirectionID,
	).Scan(&c); err != nil {
		return 0, err
	}

	return c, nil
}

func (r *DirectionsMaRepository) Selection(querys interface{}) ([]*models.DirectionMA, error) {
	q := querys.(*models.DirectionMASelection)
	arr := []*models.DirectionMA{}

	rows, err := r.store.Query(
		`
		SELECT id, direction_id, ma_id, service_type, status, created_at, updated_at
		FROM directions_ma
		WHERE direction_id=$1
		ORDER BY id DESC
		OFFSET $2
		LIMIT $3
		`,
		q.DirectionID,
		AppMath.OffsetThreshold(*q.Page, *q.Limit),
		*q.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			dma := &models.DirectionMA{}
			if err := rows.Scan(
				&dma.ID,
				&dma.DirectionID,
				&dma.MaID,
				&dma.ServiceType,
				&dma.Status,
				&dma.CreatedAt,
				&dma.UpdatedAt,
			); err != nil {
				continue
			}

			arr = append(arr, dma)
		}

		return arr, nil
	}

	return nil, AppError.ErrInvalidCondition
}
