package notification

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
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
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидация типа уведомления
	if err := r.NotificationTypeValidation(); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Получаю всех менеджеров из БД
	uArr, err := m.store.User().GetAllManagers()
	if err != nil {
		fmt.Println(err)
	}

	switch r.Type {
	case static.NTF__T__VERIFICATION:
		// Запись уведомления в очередь для всех менеджеров
		for i := 0; i < len(uArr); i++ {
			payload, err := json.Marshal(newVefificationNotify(uArr, i, r))
			if err != nil {
				fmt.Println(err)
			}

			if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
				fmt.Println(err)
			}
		}

	case static.NTF__T__EXCHANGE_ERROR:
		// Запись уведомления в очередь для всех менеджеров
		for i := 0; i < len(uArr); i++ {
			payload, err := json.Marshal(newActionCancelNotify(uArr, i, r))
			if err != nil {
				fmt.Println(err)
			}

			if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
				fmt.Println(err)
			}
		}

	case static.NTF__T__REQ_SUPPORT:
		// Запись уведомления в очередь для всех менеджеров
		for i := 0; i < len(uArr); i++ {
			payload, err := json.Marshal(newSupportReqNotify(uArr, i, r))
			if err != nil {
				fmt.Println(err)
			}

			if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
				fmt.Println(err)
			}
		}
	}

	m.responser.NewRecord(c, r, m.store.AdminPanel().Notification().Create(r))
}
