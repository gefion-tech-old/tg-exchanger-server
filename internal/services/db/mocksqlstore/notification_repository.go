package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type NotificationRepository struct {
	notification map[uint]*models.Notification
}

func (r *NotificationRepository) Create(n *models.Notification) (*models.Notification, error) {
	return nil, nil
}

func (r *NotificationRepository) Delete(n *models.Notification) (*models.Notification, error) {
	return nil, nil
}
