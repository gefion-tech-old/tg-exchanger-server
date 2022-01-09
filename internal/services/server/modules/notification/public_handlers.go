package notification

import (
	"encoding/json"
	"net/http"
	"reflect"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppTypes "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
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
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.Validation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Notification{}) {
			return
		}

		m.responser.CreateRecordResponse(c, m.store.AdminPanel().Notification(), r,
			func() error {
				// Получаю всех менеджеров из БД
				uArr, err := m.store.User().GetAllManagers()
				if err != nil {
					return err
				}

				switch r.Type {
				case AppTypes.NotifyTypeVerification:
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

				case AppTypes.NotifyTypeExchangeError:
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

				case AppTypes.NotifyTypeReqSupport:
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
