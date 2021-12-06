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

func (pr *PrivateRoutes) updateAllBotMessageHandler(c *gin.Context) {
	req := &models.BotMessage{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
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

func (pr *PrivateRoutes) createNewBotMessageHandler(c *gin.Context) {
	req := &models.BotMessage{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	// Проверяю может ли этот пользователь создавать сообщения
	if err := req.BotMessageValidation(pr.users.Managers, pr.users.Developers); err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "you cannot edit this resource",
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
