package notification

import (
	"encoding/json"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path admin/notification
	@Type PUBLIC
	@Documentation

	Создать запись в таблице `notifications`

	# TESTED
*/
func (m *ModNotification) CreateNotificationHandler(c *gin.Context) {
	r := &models.Notification{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидация типа уведомления
	if err := m.responser.Error(c, http.StatusUnprocessableEntity,
		r.NotificationTypeValidation()); err != nil {
		return
	}

	if obj := m.responser.RecordHandler(c, r,
		r.NotificationTypeValidation(),
	); obj != nil {
		m.responser.CreateRecordResponse(c, m.store.AdminPanel().Notification(), r,
			func() error {
				// Получаю всех менеджеров из БД
				uArr, err := m.store.User().GetAllManagers()
				if err != nil {
					return err
				}

				switch r.Type {
				case static.NTF__T__VERIFICATION:
					// Запись уведомления в очередь для всех менеджеров
					for i := 0; i < len(uArr); i++ {
						payload, err := json.Marshal(newVefificationNotify(uArr, i, r))
						if err != nil {
							return err
						}

						if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
							return err
						}
					}

				case static.NTF__T__EXCHANGE_ERROR:
					// Запись уведомления в очередь для всех менеджеров
					for i := 0; i < len(uArr); i++ {
						payload, err := json.Marshal(newActionCancelNotify(uArr, i, r))
						if err != nil {
							return err
						}

						if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
							return err
						}
					}

				case static.NTF__T__REQ_SUPPORT:
					// Запись уведомления в очередь для всех менеджеров
					for i := 0; i < len(uArr); i++ {
						payload, err := json.Marshal(newSupportReqNotify(uArr, i, r))
						if err != nil {
							return err
						}

						if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
							return err
						}
					}
				}

				return nil
			},
		)
	}
}
