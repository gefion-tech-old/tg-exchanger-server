package notification

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

/*
	@Method GET
	@Path admin/notifications?page=1&limit=15
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `notifications`

	# TESTED

*/
func (m *ModNotification) GetNotificationsSelectionHandler(c *gin.Context) {
	errs, _ := errgroup.WithContext(c)

	cArrN := make(chan []*models.Notification)
	cCount := make(chan *int)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Достаю из БД запрашиваемые записи
	errs.Go(func() error {
		defer close(cArrN)
		arrN, err := m.store.AdminPanel().Notification().Selection(page, limit)
		if err != nil {
			return err
		}

		cArrN <- arrN
		return nil
	})

	// Подсчет кол-ва уведомлений в таблице
	errs.Go(func() error {
		defer close(cCount)
		c, err := m.store.AdminPanel().Notification().Count()
		if err != nil {
			return err
		}

		cCount <- &c
		return nil
	})

	arrN := <-cArrN
	count := <-cCount

	if arrN == nil || count == nil {
		tools.ServErr(c, http.StatusInternalServerError, errs.Wait())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"current_page": page,
		"last_page":    math.Ceil(float64(*count) / float64(limit)),
		"data":         arrN,
	})
}

/*
	@Method PUT
	@Path admin/notification
	@Type PRIVATE
	@Documentation

	Обновить поле status записи в таблице `notifications`

	# TESTED
*/
func (m *ModNotification) UpdateNotificationStatusHandler(c *gin.Context) {
	req := &models.Notification{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if req.NotificationStatusValidation() != nil || req.NotificationTypeValidation() != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	req.ID = id

	n, err := m.store.AdminPanel().Notification().UpdateStatus(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, n)
		return
	case sql.ErrNoRows:
		tools.ServErr(c, http.StatusNotFound, errors.ErrRecordNotFound)
		return
	default:
		tools.ServErr(c, http.StatusNotFound, err)
		return
	}
}

/*
	@Method DELETE
	@Path admin/notification
	@Type PRIVATE
	@Documentation

	Удалить запись в таблице `notifications`

	# TESTED
*/
func (m *ModNotification) DeleteNotificationHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	n, err := m.store.AdminPanel().Notification().Delete(&models.Notification{ID: id})
	switch err {
	case nil:
		c.JSON(http.StatusOK, n)
		return
	case sql.ErrNoRows:
		tools.ServErr(c, http.StatusNotFound, errors.ErrRecordNotFound)
		return
	default:
		tools.ServErr(c, http.StatusInternalServerError, err)
		return
	}
}

func newSupportReqNotify(uArr []*models.User, i int, n *models.Notification) map[string]interface{} {
	return map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  uArr[i].ChatID,
			"username": uArr[i].Username,
		},
		"message": map[string]interface{}{
			"type": "confirmation_req",
			"text": fmt.Sprintf("🔵 Запрос тех. поддержки 🔵\n\n*Пользователь*: @%s", n.User.Username),
		},
		"created_at": time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
	}
}

func newVefificationNotify(uArr []*models.User, i int, n *models.Notification) map[string]interface{} {
	return map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  uArr[i].ChatID,
			"username": uArr[i].Username,
		},
		"message": map[string]interface{}{
			"type": "confirmation_req",
			"text": fmt.Sprintf("🟢 Новая заявка 🟢\n\n*Пользователь*: @%s", n.User.Username),
		},
		"created_at": time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
	}
}

func newActionCancelNotify(uArr []*models.User, i int, n *models.Notification) map[string]interface{} {
	return map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  uArr[i].ChatID,
			"username": uArr[i].Username,
		},
		"message": map[string]interface{}{
			"type": "skip_operation",
			"text": fmt.Sprintf("🔴 Отмена операции 🔴\n\n*Пользователь*: @%s", n.User.Username),
		},
		"created_at": time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
	}
}
