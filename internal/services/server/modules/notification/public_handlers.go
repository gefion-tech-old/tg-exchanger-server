package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
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
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/notifications/notifications_api.md#create

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

func newSupportReqNotify(uArr []*models.User, i int, n *models.Notification) map[string]interface{} {
	return map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  uArr[i].ChatID,
			"username": uArr[i].Username,
		},
		"message": map[string]interface{}{
			"type": AppTypes.QueueEventConfirmationRequest,
			"text": fmt.Sprintf("🔵 Запрос тех. поддержки 🔵\n\n*Пользователь*: @%s", n.User.Username),
		},
		"created_at": time.Now().UTC().Format(core.DateStandart),
	}
}

func newVefificationNotify(uArr []*models.User, i int, n *models.Notification) map[string]interface{} {
	return map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  uArr[i].ChatID,
			"username": uArr[i].Username,
		},
		"message": map[string]interface{}{
			"type": AppTypes.QueueEventConfirmationRequest,
			"text": fmt.Sprintf("🟢 Новая заявка 🟢\n\n*Пользователь*: @%s", n.User.Username),
		},
		"created_at": time.Now().UTC().Format(core.DateStandart),
	}
}

func newActionCancelNotify(uArr []*models.User, i int, n *models.Notification) map[string]interface{} {
	return map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  uArr[i].ChatID,
			"username": uArr[i].Username,
		},
		"message": map[string]interface{}{
			"type": AppTypes.QueueEventSkipOperation,
			"text": fmt.Sprintf("🔴 Отмена операции 🔴\n\n*Пользователь*: @%s", n.User.Username),
		},
		"created_at": time.Now().UTC().Format(core.DateStandart),
	}
}
