package private

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	При валидных данных в БД создается запись о новом уведомлении.
*/
func (pr *PrivateRoutes) createNotification(c *gin.Context) {
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
	uArr, err := pr.store.User().GetAllManagers()
	if err != nil {
		fmt.Println(err)
	}

	switch req.Type {
	case static.NTF__T__VERIFICATION:
		// Запись уведомления в очередь для всех менеджеров
		for i := 0; i < len(uArr); i++ {
			payload, err := json.Marshal(map[string]interface{}{
				"to": map[string]interface{}{
					"chat_id":  uArr[i].ChatID,
					"username": uArr[i].Username,
				},
				"message": map[string]interface{}{
					"type": "confirmation_req",
					"text": fmt.Sprintf("🟢 Новая заявка на врефикацию карты 🟢\n\n*Пользователь*: @%s", req.User.Username),
				},
				"created_at": time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
			})
			if err != nil {
				fmt.Println(err)
			}

			if err := pr.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
				fmt.Println(err)
			}
		}
	}

	n, err := pr.store.Manager().Notification().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, n)
		return
	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

}

/*
	@Method GET
	@Path admin/notification
	@Type PRIVATE
	@Documentation

	При валидных данных из БД достается одно запрашиваемое уведомление.
*/
func (pr *PrivateRoutes) getNotification(c *gin.Context) {

}

/*
	@Method GET
	@Path admin/notifications
	@Type PRIVATE
	@Documentation

	При валидных данных из БД достается список уведомлений
	за последний месяц.
*/
func (pr *PrivateRoutes) getAllNotifications(c *gin.Context) {}

/*
	@Method PUT
	@Path admin/notification
	@Type PRIVATE
	@Documentation

	Обновить статус уведомления.
*/
func (pr *PrivateRoutes) updateNotificationStatus(c *gin.Context) {}

/*
	@Method DELETE
	@Path admin/notification
	@Type PRIVATE
	@Documentation

	Удалить конкретное уведомление
*/
func (pr *PrivateRoutes) deleteNotification(c *gin.Context) {}
