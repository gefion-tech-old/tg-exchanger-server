package mocksqlstore

import (
	"database/sql"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
)

type NotificationRepository struct {
	notification map[int]*models.Notification
}

func (r *NotificationRepository) Create(n *models.Notification) error {
	n.ID = len(r.notification) + 1
	n.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")
	n.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

	r.notification[n.ID] = n
	return nil
}

func (r *NotificationRepository) Get(n *models.Notification) error {
	for i := 0; i < len(r.notification); i++ {
		if r.notification[i].ID == n.ID {
			r.rewrite(n.ID, n)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *NotificationRepository) Delete(n *models.Notification) error {
	if r.notification[n.ID] != nil {
		r.rewrite(n.ID, n)
		defer delete(r.notification, r.notification[n.ID].ID)
		return nil
	}

	return sql.ErrNoRows
}

func (r *NotificationRepository) UpdateStatus(n *models.Notification) error {
	if r.notification[n.ID] != nil {
		r.notification[n.ID].Status = n.Status
		r.notification[n.ID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

		r.rewrite(n.ID, n)
		return nil
	}

	return sql.ErrNoRows
}

func (r *NotificationRepository) Selection(page, limit int) ([]*models.Notification, error) {
	arr := []*models.Notification{}

	for i, v := range r.notification {
		if i > tools.OffsetThreshold(page, limit) && i <= tools.OffsetThreshold(page, limit)+limit {
			arr = append(arr, v)
		}
		i++
	}

	return arr, nil
}

func (r *NotificationRepository) Count() (int, error) {
	return len(r.notification), nil
}

func (r *NotificationRepository) rewrite(id int, to *models.Notification) {
	to.ID = r.notification[id].ID
	to.Type = r.notification[id].Type
	to.Status = r.notification[id].Status
	to.MetaData = r.notification[id].MetaData
	to.User = r.notification[id].User
	to.CreatedAt = r.notification[id].CreatedAt
	to.UpdatedAt = r.notification[id].UpdatedAt
}
