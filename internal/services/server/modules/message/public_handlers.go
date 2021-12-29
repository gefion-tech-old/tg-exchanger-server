package message

import (
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
	r := &models.BotMessage{Connector: c.Param("connector")}
	m.responser.Record(c, r, m.store.AdminPanel().BotMessages().Get(r))
}
