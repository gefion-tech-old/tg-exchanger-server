package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type MerchantAutopayoutRepository struct {
	store *sql.DB
}

func (r *MerchantAutopayoutRepository) Create(m *models.MerchantAutopayout) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO merchant_autopayout (name, service, service_type, options, status, message_id, created_by)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT name FROM merchant_autopayout WHERE name=$8)
		RETURNING id, name, service, service_type, options, status, message_id, created_by, created_at, updated_at
		`,
		m.Name,
		m.Service,
		m.ServiceType,
		m.Options,
		m.Status,
		m.MessageID,
		m.CreatedBy,
		m.Name,
	).Scan(
		&m.ID,
		&m.Name,
		&m.Service,
		&m.ServiceType,
		&m.Options,
		&m.Status,
		&m.MessageID,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *MerchantAutopayoutRepository) Update(m *models.MerchantAutopayout) error {
	if err := r.store.QueryRow(
		`
		UPDATE merchant_autopayout
		SET name=$1, service_type=$2, options=$3, status=$4, message_id=$5, updated_at=$6
		WHERE id=$7
		RETURNING id, name, service, service_type, options, status, message_id, created_by, created_at, updated_at
		`,
		m.Name,
		m.ServiceType,
		m.Options,
		m.Status,
		m.MessageID,
		time.Now().UTC().Format(core.DateStandart),
		m.ID,
	).Scan(
		&m.ID,
		&m.Name,
		&m.Service,
		&m.ServiceType,
		&m.Options,
		&m.Status,
		&m.MessageID,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *MerchantAutopayoutRepository) Get(m *models.MerchantAutopayout) error {
	if err := r.store.QueryRow(
		`
		SELECT id, name, service, service_type, options, status, message_id, created_by, created_at, updated_at
		FROM merchant_autopayout
		WHERE id=$1
		`,
		m.ID,
	).Scan(
		&m.ID,
		&m.Name,
		&m.Service,
		&m.ServiceType,
		&m.Options,
		&m.Status,
		&m.MessageID,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *MerchantAutopayoutRepository) Delete(m *models.MerchantAutopayout) error {
	if err := r.store.QueryRow(
		`
		DELETE FROM merchant_autopayout
		WHERE id=$1
		RETURNING id, name, service, service_type, options, status, message_id, created_by, created_at, updated_at
		`,
		m.ID,
	).Scan(
		&m.ID,
		&m.Name,
		&m.Service,
		&m.ServiceType,
		&m.Options,
		&m.Status,
		&m.MessageID,
		&m.CreatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *MerchantAutopayoutRepository) Selection(querys interface{}) ([]*models.MerchantAutopayout, error) {
	q := querys.(*models.MerchantAutopayoutSelection)
	arr := []*models.MerchantAutopayout{}

	sb := fmt.Sprintf(`
		SELECT id, name, service, service_type, options, status, message_id, created_by, created_at, updated_at
		FROM merchant_autopayout
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
			r := &models.MerchantAutopayout{}
			if err := rows.Scan(
				&r.ID,
				&r.Name,
				&r.Service,
				&r.ServiceType,
				&r.Options,
				&r.Status,
				&r.MessageID,
				&r.CreatedBy,
				&r.CreatedAt,
				&r.UpdatedAt,
			); err != nil {
				continue
			}

			arr = append(arr, r)
		}

		return arr, nil
	}

	return nil, AppError.ErrInvalidCondition
}

func (r *MerchantAutopayoutRepository) Count(querys interface{}) (int, error) {
	q := querys.(*models.MerchantAutopayoutSelection)
	var c int

	sb := fmt.Sprintf(`
		SELECT count(*)
		FROM merchant_autopayout
		WHERE %s
	`,
		strings.Join(r.queryGeneration(q), " AND "),
	)

	if err := r.store.QueryRow(sb).Scan(
		&c,
	); err != nil {
		return 0, err
	}

	return c, nil
}

func (r *MerchantAutopayoutRepository) queryGeneration(q *models.MerchantAutopayoutSelection) []string {
	var conditions []string

	if q.Service[0] == "" || len(q.Service) == 0 {
		q.Service = q.Service[:len(q.Service)-1]
		q.Service = append(q.Service,
			AppType.MerchantAutoPayoutWhitebit,
			AppType.MerchantAutoPayoutMine,
		)

		conditions = append(conditions, fmt.Sprintf("service='%s' OR service='%s'",
			AppType.MerchantAutoPayoutWhitebit,
			AppType.MerchantAutoPayoutMine,
		))
	} else {
		conditions = append(conditions, fmt.Sprintf("service='%s'", q.Service[0]))
	}

	return conditions
}
