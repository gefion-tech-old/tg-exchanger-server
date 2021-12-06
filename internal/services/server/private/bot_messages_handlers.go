package private

import (
	"database/sql"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (pr *PrivateRoutes) deleteBotMessageHandler(c *gin.Context) {
	req := &models.BotMessage{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	msg, err := pr.store.Manager().BotMessages().Delete(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, msg)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "message with current connector is not found",
		})
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
	@Path admin/message?connector=<connector_name>
	@Type PRIVATE
	@Documentation https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#get-all

	При валидных данных обновляет запись конкретного сообщения в БД.

	# TESTED
*/
func (pr *PrivateRoutes) updateAllBotMessageHandler(c *gin.Context) {
	req := &models.BotMessage{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	if err := req.BotMessageValidation(pr.users.Managers, pr.users.Developers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg, err := pr.store.Manager().BotMessages().Update(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, msg)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "message with current connector is not found",
		})
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
	@Path admin/message?connector=<connector_name>
	@Type PUBLIC
	@Documentation https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#get-all

	При валидных данных возвращается запрашиваемое сообщение.

	# TESTED
*/
func (pr *PrivateRoutes) getAllBotMessageHandler(c *gin.Context) {
	msgs, err := pr.store.Manager().BotMessages().GetAll()
	switch err {
	case nil:
		c.JSON(http.StatusOK, msgs)
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
	@Path admin/message?connector=<connector_name>
	@Type PUBLIC
	@Documentation https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#get

	При валидных данных возвращается запрашиваемое сообщение.

	# TESTED
*/
func (pr *PrivateRoutes) getBotMessageHandler(c *gin.Context) {
	if c.Query("connector") == "" {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	msg, err := pr.store.Manager().BotMessages().Get(&models.BotMessage{Connector: c.Query("connector")})
	switch err {
	case nil:
		c.JSON(http.StatusOK, msg)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "message with current connector is not found",
		})
		return

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

}

/*
	@Method POST
	@Path admin/message
	@Type PRIVATE
	@Documentation https://github.com/gefion-tech/tg-exchanger-server/blob/main/docs/admin__bot_messages.md#create

	При валидных данных создаем в БД запись нового сообщения
	и возвращаю созданное сообщение.

	# TESTED
*/
func (pr *PrivateRoutes) createNewBotMessageHandler(c *gin.Context) {
	req := &models.BotMessage{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	if err := req.BotMessageValidation(pr.users.Managers, pr.users.Developers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg, err := pr.store.Manager().BotMessages().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, msg)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusConflict, gin.H{
			"error": "message with current connector already created",
		})
		return

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}
