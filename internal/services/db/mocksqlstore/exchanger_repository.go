package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
)

type ExchangerRepository struct {
	exchangers map[int]*models.Exchanger
}

func (r *ExchangerRepository) Create(e *models.Exchanger) error {
	e.ID = len(r.exchangers) + 1
	e.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	e.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.exchangers[e.ID] = e
	return nil
}

func (r *ExchangerRepository) Update(e *models.Exchanger) error {
	if r.exchangers[e.ID] != nil {
		r.exchangers[e.ID].Name = e.Name
		r.exchangers[e.ID].UrlToParse = e.UrlToParse
		r.exchangers[e.ID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
		return nil

	}

	return sql.ErrNoRows
}

func (r *ExchangerRepository) GetByName(e *models.Exchanger) error {
	for _, ex := range r.exchangers {
		if ex.Name == e.Name {
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *ExchangerRepository) Delete(e *models.Exchanger) error {
	if r.exchangers[e.ID] != nil {
		defer delete(r.exchangers, r.exchangers[e.ID].ID)
		return nil
	}

	return sql.ErrNoRows
}

func (r *ExchangerRepository) Count() (int, error) {
	return len(r.exchangers), nil
}

func (r *ExchangerRepository) Selection(page, limit int) ([]*models.Exchanger, error) {
	arr := []*models.Exchanger{}

	for i, v := range r.exchangers {
		if i > tools.OffsetThreshold(page, limit) && i <= tools.OffsetThreshold(page, limit)+limit {
			arr = append(arr, v)
		}
		i++
	}

	return arr, nil
}

func (r *ExchangerRepository) GetSlice(limit int) ([]*models.Exchanger, error) {
	eArr := []*models.Exchanger{}

	for i := 0; i < limit; i++ {
		eArr = append(eArr, r.exchangers[i])
	}

	return eArr, nil
}
