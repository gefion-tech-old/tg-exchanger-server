package notification

import (
	"net/http"
	"reflect"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method GET
	@Path admin/notifications/check
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/notifications/notifications_api.md#count-new-notifications

	Получить кол-во новых уведомлений `notifications`

	# TESTED
*/
func (m *ModNotification) NewNotificationsCheckHandler(c *gin.Context) {
	count, err := m.store.AdminPanel().Notification().CheckNew()
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"new_notifications": count,
	})
}

/*
	@Method GET
	@Path admin/notifications?page=1&limit=15
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/notifications/notifications_api.md#selection

	Получение лимитированного объема записей из таблицы `notifications`

	# TESTED
*/
func (m *ModNotification) GetNotificationsSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c,
		m.store.AdminPanel().Notification(),
		&models.NotificationSelection{},
	)
}

/*
	@Method PUT
	@Path admin/notification
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/notifications/notifications_api.md#update

	Обновить поле status записи в таблице `notifications`

	# TESTED
*/
func (m *ModNotification) UpdateNotificationStatusHandler(c *gin.Context) {
	r := &models.Notification{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.Validation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Notification{}) {
			return
		}

		m.responser.UpdateRecordResponse(c, m.store.AdminPanel().Notification(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	@Method DELETE
	@Path admin/notification
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/notifications/notifications_api.md#delete

	Удалить запись в таблице `notifications`

	# TESTED
*/
func (m *ModNotification) DeleteNotificationHandler(c *gin.Context) {
	if obj := m.responser.RecordHandler(c, &models.Notification{}); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Notification{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.store.AdminPanel().Notification(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}
