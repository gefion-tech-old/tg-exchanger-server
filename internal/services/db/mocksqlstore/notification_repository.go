package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type NotificationRepository struct {
	notification map[int]*models.Notification
}

func (r *NotificationRepository) Create(n *models.Notification) (*models.Notification, error) {
	n.ID = len(r.notification) + 1
	n.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	n.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.notification[n.ID] = n
	return r.notification[n.ID], nil
}

func (r *NotificationRepository) Get(n *models.Notification) (*models.Notification, error) {
	for i := 0; i < len(r.notification); i++ {
		if r.notification[i].ID == n.ID {
			return r.notification[i], nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *NotificationRepository) Delete(n *models.Notification) (*models.Notification, error) {
	if r.notification[n.ID] != nil {
		defer delete(r.notification, r.notification[n.ID].ID)
		return r.notification[n.ID], nil
	}

	return nil, sql.ErrNoRows
}

func (r *NotificationRepository) UpdateStatus(n *models.Notification) (*models.Notification, error) {
	if r.notification[n.ID] != nil {
		r.notification[n.ID].Status = n.Status
		r.notification[n.ID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
		return r.notification[n.ID], nil
	}

	return nil, sql.ErrNoRows
}

func (r *NotificationRepository) GetSlice(limit int) ([]*models.Notification, error) {
	nArr := []*models.Notification{}

	for i := 0; i < limit; i++ {
		nArr = append(nArr, r.notification[i])
	}

	return nArr, nil
}

func (r *NotificationRepository) Count() (int, error) {
	return len(r.notification), nil
}
