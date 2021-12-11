package private

import (
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	// Валидация типа уведомления
	if err := req.NotificationTypeValidation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	n, err := pr.store.Manager().Notification().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, n)
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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
func (pr *PrivateRoutes) getNotification(c *gin.Context) {}

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
