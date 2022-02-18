package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type DirectionsRepository struct {
	store *sql.DB

	directionsMaRepository *DirectionsMaRepository
}

/*
	==========================================================================================
	КОНСТРУКТОРЫ ВЛОЖЕННЫХ СТРУКТУР
	==========================================================================================
*/

func (r *DirectionsRepository) Ma() db.DirectionsMaRepository {
	if r.directionsMaRepository != nil {
		return r.directionsMaRepository
	}

	r.directionsMaRepository = &DirectionsMaRepository{
		store: r.store,
	}

	return r.directionsMaRepository
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ
	==========================================================================================
*/
func (r *DirectionsRepository) Create(d *models.Direction) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO exchange_directions(exchange_from, exchange_to, course_correction, address_verification, status, created_by)
		SELECT $1, $2, $3, $4, $5, $6
		RETURNING id, exchange_from, exchange_to, course_correction, address_verification, status, created_by, created_at, updated_at
		`,
		d.ExchangeFrom,
		d.ExchangeTo,
		d.CourseCorrection,
		d.AddressVerification,
		d.Status,
		d.CreatedBy,
	).Scan(
		&d.ID,
		&d.ExchangeFrom,
		&d.ExchangeTo,
		&d.CourseCorrection,
		&d.AddressVerification,
		&d.Status,
		&d.CreatedBy,
		&d.CreatedAt,
		&d.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsRepository) Update(d *models.Direction) error {
	if err := r.store.QueryRow(
		`
		UPDATE exchange_directions
		SET exchange_from=$1, exchange_to=$2, course_correction=$3, address_verification=$4, status=$5, updated_at=$6
		WHERE id=$7
		RETURNING id, exchange_from, exchange_to, course_correction, address_verification, status, created_by, created_at, updated_at
		`,
		d.ExchangeFrom,
		d.ExchangeTo,
		d.CourseCorrection,
		d.AddressVerification,
		d.Status,
		time.Now().UTC().Format(core.DateStandart),
		d.ID,
	).Scan(
		&d.ID,
		&d.ExchangeFrom,
		&d.ExchangeTo,
		&d.CourseCorrection,
		&d.AddressVerification,
		&d.Status,
		&d.CreatedBy,
		&d.CreatedAt,
		&d.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsRepository) Delete(d *models.Direction) error {
	if err := r.store.QueryRow(
		`
		DELETE FROM exchange_directions
		WHERE id=$1
		RETURNING id, exchange_from, exchange_to, course_correction, address_verification, status, created_by, created_at, updated_at
		`,
		d.ID,
	).Scan(
		&d.ID,
		&d.ExchangeFrom,
		&d.ExchangeTo,
		&d.CourseCorrection,
		&d.AddressVerification,
		&d.Status,
		&d.CreatedBy,
		&d.CreatedAt,
		&d.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsRepository) Get(d *models.Direction) error {
	if err := r.store.QueryRow(
		`
		SELECT id, exchange_from, exchange_to, course_correction, address_verification, status, created_by, created_at, updated_at
		FROM exchange_directions
		WHERE id=$1
		`,
		d.ID,
	).Scan(
		&d.ID,
		&d.ExchangeFrom,
		&d.ExchangeTo,
		&d.CourseCorrection,
		&d.AddressVerification,
		&d.Status,
		&d.CreatedBy,
		&d.CreatedAt,
		&d.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *DirectionsRepository) Count(querys interface{}) (int, error) {
	q := querys.(*models.DirectionSelection)
	var c int

	sb := fmt.Sprintf(`
	SELECT count(*)
	FROM exchange_directions
	WHERE %s
	`,
		strings.Join(r.queryGeneration(q), " AND "),
	)

	if err := r.store.QueryRow(sb).Scan(&c); err != nil {
		return 0, err
	}

	return c, nil
}

func (r *DirectionsRepository) Selection(querys interface{}) ([]*models.Direction, error) {
	q := querys.(*models.DirectionSelection)
	arr := []*models.Direction{}

	sb := fmt.Sprintf(`
		SELECT id, exchange_from, exchange_to, course_correction, address_verification, status, created_by, created_at, updated_at
		FROM exchange_directions
		WHERE %s
		ORDER BY id DESC
		OFFSET %d
		LIMIT %d
	`,
		strings.Join(r.queryGeneration(q), " AND "),
		AppMath.OffsetThreshold(*q.Page, *q.Limit),
		*q.Limit,
	)

	rows, err := r.store.Query(sb)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			d := &models.Direction{}
			if err := rows.Scan(
				&d.ID,
				&d.ExchangeFrom,
				&d.ExchangeTo,
				&d.CourseCorrection,
				&d.AddressVerification,
				&d.Status,
				&d.CreatedBy,
				&d.CreatedAt,
				&d.UpdatedAt,
			); err != nil {
				continue
			}

			arr = append(arr, d)
		}

		return arr, nil
	}

	return nil, AppError.ErrInvalidCondition
}

func (r *DirectionsRepository) queryGeneration(q *models.DirectionSelection) []string {
	conditions := []string{}

	if q.Status != nil {
		switch *q.Status {
		case false:
			conditions = append(conditions, "status=FALSE")
		case true:
			conditions = append(conditions, "status=TRUE")
		}
	}

	return conditions
}
