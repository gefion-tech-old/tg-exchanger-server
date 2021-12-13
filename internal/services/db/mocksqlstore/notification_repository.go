package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type NotificationRepository struct {
	notification map[uint]*models.Notification
}

func (r *NotificationRepository) Create(n *models.Notification) (*models.Notification, error) {
	n.ID = uint(len(r.notification) + 1)
	n.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	n.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.notification[n.ID] = n
	return r.notification[n.ID], nil
}

func (r *NotificationRepository) Get(n *models.Notification) (*models.Notification, error) {
	for i := 0; i < len(r.notification); i++ {
		if r.notification[uint(i)].ID == n.ID {
			return r.notification[uint(i)], nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *NotificationRepository) Delete(n *models.Notification) (*models.Notification, error) {
	if r.notification[uint(n.ID)] != nil {
		defer delete(r.notification, r.notification[uint(n.ID)].ID)
		return r.notification[uint(n.ID)], nil
	}

	return nil, sql.ErrNoRows
}

func (r *NotificationRepository) UpdateStatus(n *models.Notification) (*models.Notification, error) {
	if r.notification[uint(n.ID)] != nil {
		r.notification[uint(n.ID)].Status = n.Status
		r.notification[uint(n.ID)].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
		return r.notification[uint(n.ID)], nil
	}

	return nil, sql.ErrNoRows
}

func (r *NotificationRepository) GetSlice(limit int) ([]*models.Notification, error) {
	nArr := []*models.Notification{}

	for i := 0; i < limit; i++ {
		nArr = append(nArr, r.notification[uint(i)])
	}

	return nArr, nil
}

func (r *NotificationRepository) Count() (int, error) {
	return len(r.notification), nil
}
