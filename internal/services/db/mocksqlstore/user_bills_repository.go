package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserBillsRepository struct {
	bills map[uint]*models.Bill
}

func (r *UserBillsRepository) Create(b *models.Bill) (*models.Bill, error) {
	b.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	b.ID = uint(len(r.bills) + 1)

	r.bills[b.ID] = b
	return r.bills[b.ID], nil
}

func (r *UserBillsRepository) Delete(b *models.Bill) (*models.Bill, error) {
	for _, bill := range r.bills {
		if bill.Bill == b.Bill && bill.ChatID == b.ChatID {
			delete(r.bills, bill.ID)

			b.ID = bill.ID
			b.CreatedAt = bill.CreatedAt
			return b, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *UserBillsRepository) All(chatID int64) ([]*models.Bill, error) {
	arr := []*models.Bill{}
	for _, bill := range r.bills {
		if bill.ChatID == chatID {
			arr = append(arr, bill)
		}
	}
	return arr, nil
}
