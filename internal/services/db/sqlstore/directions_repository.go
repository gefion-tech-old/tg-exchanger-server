package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type DirectionsRepository struct {
	store *sql.DB
}

func (r *DirectionsRepository) Create(d *models.Direction) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO exchange_directions (exchange_from, exchange_to, course_correction, address_verification, direction_status)
		SELECT $1, $2, $3, $4, $5
		RETURNING id, exchange_from, exchange_to, course_correction, address_verification, direction_status, created_at, updated_at
		`,
		d.ExchangeFrom,
		d.ExchangeTo,
		d.CourseCorrection,
		d.AddressVerification,
		d.DirectionStatus,
	).Scan(
		&d.ID,
		&d.ExchangeFrom,
		&d.ExchangeTo,
		&d.CourseCorrection,
		&d.AddressVerification,
		&d.DirectionStatus,
		&d.CreatedAt,
		&d.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (r *DirectionsRepository) Update(d *models.Direction) error {
	return nil
}

func (r *DirectionsRepository) Delete(d *models.Direction) error {
	return nil
}

func (r *DirectionsRepository) Get(d *models.Direction) error {
	return nil
}
