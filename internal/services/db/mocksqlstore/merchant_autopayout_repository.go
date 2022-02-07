package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type MerchantAutopayoutRepository struct {
	ma map[int]*models.MerchantAutopayout
}

func (r *MerchantAutopayoutRepository) Create(m *models.MerchantAutopayout) error {
	m.ID = len(r.ma) + 1
	m.CreatedAt = time.Now().UTC().Format(core.DateStandart)
	m.UpdatedAt = time.Now().UTC().Format(core.DateStandart)

	r.ma[m.ID] = m
	return nil
}

func (r *MerchantAutopayoutRepository) Update(m *models.MerchantAutopayout) error {
	for _, v := range r.ma {
		if v.ID == m.ID {
			v.Name = m.Name
			v.ServiceType = m.ServiceType
			v.Options = m.Options
			v.Status = m.Status
			v.MessageID = m.MessageID
			v.UpdatedAt = time.Now().UTC().Format(core.DateStandart)

			r.rewrite(m.ID, m)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *MerchantAutopayoutRepository) Get(m *models.MerchantAutopayout) error {
	for _, v := range r.ma {
		if v.ID == m.ID {
			r.rewrite(m.ID, m)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *MerchantAutopayoutRepository) Delete(m *models.MerchantAutopayout) error {
	for _, v := range r.ma {
		if v.ID == m.ID {
			r.rewrite(m.ID, m)
			defer delete(r.ma, m.ID)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *MerchantAutopayoutRepository) Count(querys interface{}) (int, error) {
	return len(r.ma), nil
}

func (r *MerchantAutopayoutRepository) Selection(querys interface{}) ([]*models.MerchantAutopayout, error) {
	arr := []*models.MerchantAutopayout{}

	for i := 0; i < len(r.ma); i++ {
		arr = append(arr, r.ma[i])
	}

	return arr, nil
}

func (r *MerchantAutopayoutRepository) GetFistIfActive(service string) (*models.MerchantAutopayout, error) {
	return nil, nil
}

func (r *MerchantAutopayoutRepository) GetAllByServiceType(serviceType int, status bool) ([]*models.MerchantAutopayout, error) {
	return nil, nil
}

/*
	==========================================================================================
	ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
	==========================================================================================
*/

func (r *MerchantAutopayoutRepository) rewrite(id int, to *models.MerchantAutopayout) {
	to.ID = r.ma[id].ID
	to.Name = r.ma[id].Name
	to.Service = r.ma[id].Service
	to.ServiceType = r.ma[id].ServiceType
	to.Options = r.ma[id].Options
	to.Status = r.ma[id].Status
	to.MessageID = r.ma[id].MessageID
	to.CreatedBy = r.ma[id].CreatedBy
	to.CreatedAt = r.ma[id].CreatedAt
	to.UpdatedAt = r.ma[id].UpdatedAt
}
