package message

import (
	"database/sql"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method GET
	@Path admin/message/:connector
	@Type PUBLIC
	@Documentation

	Получить запись из таблицы `bot_messages`

	# TESTED
*/
func (m *ModMessage) GetMessageHandler(c *gin.Context) {
	msg, err := m.store.AdminPanel().BotMessages().Get(&models.BotMessage{Connector: c.Param("connector")})
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
