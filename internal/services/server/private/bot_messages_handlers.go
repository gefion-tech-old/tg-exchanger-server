package private

import (
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

}
