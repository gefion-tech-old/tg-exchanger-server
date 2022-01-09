package mocksqlstore

import (
	"database/sql"
	"time"

	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type NotificationRepository struct {
	notification map[int]*models.Notification
}

func (r *NotificationRepository) CheckNew() (int, error) {
	var c int

	for _, n := range r.notification {
		if n.Status == 1 {
			c++
		}
	}

	return c, nil
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

func (r *NotificationRepository) Update(n *models.Notification) error {
	if r.notification[n.ID] != nil {
		r.notification[n.ID].Status = n.Status
		r.notification[n.ID].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.00000000")

		r.rewrite(n.ID, n)
		return nil
	}

	return sql.ErrNoRows
}

func (r *NotificationRepository) Selection(querys interface{}) ([]*models.Notification, error) {
	q := querys.(*models.NotificationSelection)
	arr := []*models.Notification{}

	for i, v := range r.notification {
		if i > AppMath.OffsetThreshold(q.Page, q.Limit) && i <= AppMath.OffsetThreshold(q.Page, q.Limit)+q.Limit {
			arr = append(arr, v)
		}
		i++
	}

	return arr, nil
}

func (r *NotificationRepository) Count(querys interface{}) (int, error) {
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
