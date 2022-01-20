package message

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method GET
	@Path admin/message/:connector
	@Type PUBLIC
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/messages/messages_api.md#get

	Получить запись из таблицы `bot_messages`

	# TESTED
*/
func (m *ModMessage) GetMessageHandler(c *gin.Context) {
	r := &models.BotMessage{Connector: c.Param("connector")}
	m.responser.RecordResponse(c, r, m.store.AdminPanel().BotMessages().Get(r))
}
