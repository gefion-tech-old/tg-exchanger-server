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
	req := &models.Notification{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидация типа уведомления
	if err := req.NotificationTypeValidation(); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Получаю всех менеджеров из БД
	uArr, err := m.store.User().GetAllManagers()
	if err != nil {
		fmt.Println(err)
	}

	switch req.Type {
	case static.NTF__T__VERIFICATION:
		// Запись уведомления в очередь для всех менеджеров
		for i := 0; i < len(uArr); i++ {
			payload, err := json.Marshal(newVefificationNotify(uArr, i, req))
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
			payload, err := json.Marshal(newActionCancelNotify(uArr, i, req))
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
			payload, err := json.Marshal(newSupportReqNotify(uArr, i, req))
			if err != nil {
				fmt.Println(err)
			}

			if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
				fmt.Println(err)
			}
		}
	}

	n, err := m.store.AdminPanel().Notification().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, n)
		return
	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

}
