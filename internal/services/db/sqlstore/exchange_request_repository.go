package sqlstore

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type ExchangeRequestRepository struct {
	store *sql.DB
}

func (r *ExchangeRequestRepository) Create(er *models.ExchangeRequest) error {
	if err := r.store.QueryRow(
		`
		INSERT INTO request(request_status, exchange_from, exchange_to, course, address, client_address, expected_amount, created_by_username, created_by_chat_id)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9
		RETURNING id, request_status, exchange_from, exchange_to, course, address, client_address, expected_amount, transferred_amount, transaction_hash, created_by_username, created_by_chat_id, created_at, updated_at
		`,
		er.Status,
		er.ExchangeFrom,
		er.ExchangeTo,
		er.Course,
		er.Address,
		er.ClientAddress,
		er.ExpectedAmount,
		er.CreatedBy.Username,
		er.CreatedBy.ChatID,
	).Scan(
		&er.ID,
		&er.Status,
		&er.ExchangeFrom,
		&er.ExchangeTo,
		&er.Course,
		&er.Address,
		&er.ClientAddress,
		&er.ExpectedAmount,
		&er.TransferredAmount,
		&er.TransactionHash,
		&er.CreatedBy.Username,
		&er.CreatedBy.ChatID,
		&er.CreatedAt,
		&er.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *ExchangeRequestRepository) Get(er *models.ExchangeRequest) error {
	if err := r.store.QueryRow(
		`
		SELECT id, request_status, exchange_from, exchange_to, course, address, client_address, expected_amount, transferred_amount, transaction_hash, created_by_username, created_by_chat_id, created_at, updated_at
		FROM request
		WHERE id=$1
		`,
		er.ID,
	).Scan(
		&er.ID,
		&er.Status,
		&er.ExchangeFrom,
		&er.ExchangeTo,
		&er.Course,
		&er.Address,
		&er.ClientAddress,
		&er.ExpectedAmount,
		&er.TransferredAmount,
		&er.TransactionHash,
		&er.CreatedBy.Username,
		&er.CreatedBy.ChatID,
		&er.CreatedAt,
		&er.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *ExchangeRequestRepository) Update(er *models.ExchangeRequest) error {
	if err := r.store.QueryRow(
		`
		UPDATE request
		SET request_status=$1, transferred_amount=$2, transaction_hash=$3, updated_at=$4
		WHERE id=$5
		RETURNING id, request_status, exchange_from, exchange_to, course, address, client_address, expected_amount, transferred_amount, transaction_hash, created_by_username, created_by_chat_id, created_at, updated_at
		`,
		er.Status,
		er.TransferredAmount,
		*er.TransactionHash,
		time.Now().UTC().Format(core.DateStandart),
		er.ID,
	).Scan(
		&er.ID,
		&er.Status,
		&er.ExchangeFrom,
		&er.ExchangeTo,
		&er.Course,
		&er.Address,
		&er.ClientAddress,
		&er.ExpectedAmount,
		&er.TransferredAmount,
		&er.TransactionHash,
		&er.CreatedBy.Username,
		&er.CreatedBy.ChatID,
		&er.CreatedAt,
		&er.UpdatedAt,
	); err != nil {

		return err
	}

	return nil
}

func (r *ExchangeRequestRepository) Delete(er *models.ExchangeRequest) error {
	if err := r.store.QueryRow(
		`
		DELETE FROM request
		WHERE id=$1
		RETURNING id, request_status, exchange_from, exchange_to, course, address, client_address, expected_amount, transferred_amount, transaction_hash, created_by_username, created_by_chat_id, created_at, updated_at
		`,
		er.ID,
	).Scan(
		&er.ID,
		&er.Status,
		&er.ExchangeFrom,
		&er.ExchangeTo,
		&er.Course,
		&er.Address,
		&er.ClientAddress,
		&er.ExpectedAmount,
		&er.TransferredAmount,
		&er.TransactionHash,
		&er.CreatedBy.Username,
		&er.CreatedBy.ChatID,
		&er.CreatedAt,
		&er.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *ExchangeRequestRepository) Count(querys interface{}) (int, error) {
	q := querys.(*models.ExchangeRequestSelection)
	var c int

	sb := fmt.Sprintf(`
		SELECT count(*)
		FROM request
		%s
	`,
		r.queryGeneration(q),
	)

	if err := r.store.QueryRow(sb).Scan(&c); err != nil {
		return 0, err
	}

	return c, nil
}

func (r *ExchangeRequestRepository) Selection(querys interface{}) ([]*models.ExchangeRequest, error) {
	q := querys.(*models.ExchangeRequestSelection)
	arr := []*models.ExchangeRequest{}

	sb := fmt.Sprintf(`
		SELECT id, request_status, exchange_from, exchange_to, course, address, client_address, expected_amount, transferred_amount, transaction_hash, created_by_username, created_by_chat_id, created_at, updated_at
		FROM request
		%s
		ORDER BY id DESC
		OFFSET %d
		LIMIT %d
	`,
		r.queryGeneration(q),
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
			er := &models.ExchangeRequest{}
			if err := rows.Scan(
				&er.ID,
				&er.Status,
				&er.ExchangeFrom,
				&er.ExchangeTo,
				&er.Course,
				&er.Address,
				&er.ClientAddress,
				&er.ExpectedAmount,
				&er.TransferredAmount,
				&er.TransactionHash,
				&er.CreatedBy.Username,
				&er.CreatedBy.ChatID,
				&er.CreatedAt,
				&er.UpdatedAt,
			); err != nil {
				continue
			}

			arr = append(arr, er)
		}

		return arr, nil
	}

	return nil, AppError.ErrInvalidCondition
}

func (r *ExchangeRequestRepository) GetAllByStatus(status ...AppType.ExchangeRequestStatus) ([]*models.ExchangeRequest, error) {
	arr := []*models.ExchangeRequest{}

	sb := fmt.Sprintf(`
	SELECT id, request_status, exchange_from, exchange_to, course, address, client_address, expected_amount, transferred_amount, transaction_hash, created_by_username, created_by_chat_id, created_at, updated_at
	FROM request
	WHERE %s
	ORDER BY id DESC
	`,
		r.statusQueryGeneration(status...),
	)

	rows, err := r.store.Query(sb)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			er := &models.ExchangeRequest{}
			if err := rows.Scan(
				&er.ID,
				&er.Status,
				&er.ExchangeFrom,
				&er.ExchangeTo,
				&er.Course,
				&er.Address,
				&er.ClientAddress,
				&er.ExpectedAmount,
				&er.TransferredAmount,
				&er.TransactionHash,
				&er.CreatedBy.Username,
				&er.CreatedBy.ChatID,
				&er.CreatedAt,
				&er.UpdatedAt,
			); err != nil {
				continue
			}

			arr = append(arr, er)
		}

		return arr, nil
	}

	return arr, nil
}

/*
	==========================================================================================
	ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
	==========================================================================================
*/

func (r *ExchangeRequestRepository) statusQueryGeneration(status ...AppType.ExchangeRequestStatus) string {
	arr := []string{}
	for _, s := range status {
		arr = append(arr, fmt.Sprintf("request_status=%d", s))
	}

	return strings.Join(arr, " OR ")
}

func (r *ExchangeRequestRepository) queryGeneration(q *models.ExchangeRequestSelection) string {
	if q.Status != 0 {
		return "WHERE request_status=" + strconv.Itoa(int(q.Status))
	}

	return ""
}
