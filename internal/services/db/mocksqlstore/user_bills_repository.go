package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type UserBillsRepository struct {
	bills map[int]*models.Bill
}

func (r *UserBillsRepository) Create(b *models.Bill) error {
	b.ID = len(r.bills) + 1
	b.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.bills[b.ID] = b
	return nil
}

func (r *UserBillsRepository) FindById(b *models.Bill) error {
	for _, v := range r.bills {
		if v.ID != b.ID {

			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *UserBillsRepository) Delete(b *models.Bill) error {
	for i, bill := range r.bills {
		if bill.ID == b.ID && bill.ChatID == b.ChatID {
			r.rewrite(i, b)
			delete(r.bills, bill.ID)
			return nil
		}
	}

	return sql.ErrNoRows
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

func (r *UserBillsRepository) rewrite(id int, to *models.Bill) {
	to.ID = r.bills[id].ID
	to.Bill = r.bills[id].Bill
	to.ChatID = r.bills[id].ChatID
	to.CreatedAt = r.bills[id].CreatedAt
}
