package sqlstore

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
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
func (r *LoggerRepository) Count(querys interface{}) (int, error) {
	q := querys.(*models.LogRecordSelection)
	var c int

	sb := fmt.Sprintf(`
		SELECT count(*)
		FROM logs
		WHERE %s
	`,
		strings.Join(r.queryGeneration(q), " AND "),
	)

	if err := r.store.QueryRow(sb, q.DateFrom, q.DateTo).Scan(
		&c,
	); err != nil {
		return 0, err
	}

	return c, nil
}

// # TESTED
func (r *LoggerRepository) Selection(querys interface{}) ([]*models.LogRecord, error) {
	q := querys.(*models.LogRecordSelection)
	arr := []*models.LogRecord{}

	sb := fmt.Sprintf(`
		SELECT id, username, info, service, module, created_at
		FROM logs
		WHERE %s
		ORDER BY id DESC
		OFFSET %d
		LIMIT %d
	`,
		strings.Join(r.queryGeneration(q), " AND "),
		tools.OffsetThreshold(*q.Page, *q.Limit),
		q.Limit,
	)

	rows, err := r.store.Query(sb, q.DateFrom, q.DateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return nil, AppError.ErrInvalidCondition
}

// # TESTED
func (r *LoggerRepository) DeleteSelection(querys interface{}) ([]*models.LogRecord, error) {
	q := querys.(*models.LogRecordSelection)
	arr := []*models.LogRecord{}

	if q.DateFrom == "" {
		q.DateFrom = "2020-01-01T00:00:00.00000000"
	}

	if q.DateTo == "" {
		q.DateTo = time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000")

	}

	rows, err := r.store.Query(
		`
		DELETE FROM logs
		WHERE created_at >= $1 AND created_at < $2	
		RETURNING id, username, info, service, module, created_at
		`,
		q.DateFrom,
		q.DateTo,
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

/*
	==========================================================================================
	ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
	==========================================================================================
*/

func (r *LoggerRepository) queryGeneration(q *models.LogRecordSelection) []string {
	if q.DateFrom == "" {
		q.DateFrom = "2020-01-01T00:00:00.00000000"
	}

	if q.DateTo == "" {
		q.DateTo = time.Now().Add(27 * time.Hour).UTC().Format("2006-01-02T15:04:05.00000000")
	}

	conditions := []string{"created_at >= $1 AND created_at < $2"}

	if q.Username != "" {
		conditions = append(conditions, fmt.Sprintf("username='%s'", q.Username))
	}

	if len(q.Service) > 0 {
		if q.Service[0] == "" {
			q.Service = q.Service[:len(q.Service)-1]

			q.Service = append(q.Service,
				strconv.Itoa(static.L__ADMIN),
				strconv.Itoa(static.L__BOT),
				strconv.Itoa(static.L__SERVER),
			)
		}

		conditions = append(conditions, fmt.Sprintf("service::int IN(%s)", strings.Join(q.Service, ", ")))
	}

	return conditions
}
