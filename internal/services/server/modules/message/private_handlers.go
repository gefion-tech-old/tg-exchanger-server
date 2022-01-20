package message

import (
	"net/http"
	"reflect"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method DELETE
	@Path admin/message/:id
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/messages/messages_api.md#delete

	Удалить запись в таблице `bot_messages`

	# TESTED
*/
func (m *ModMessage) DeleteBotMessageHandler(c *gin.Context) {
	if obj := m.responser.RecordHandler(c, &models.BotMessage{}); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.BotMessage{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.store.AdminPanel().BotMessages(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	@Method PUT
	@Path admin/message/:id
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/messages/messages_api.md#update

	Обновить запись в таблице `bot_messages`

	# TESTED
*/
func (m *ModMessage) UpdateBotMessageHandler(c *gin.Context) {
	r := &models.BotMessage{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.Validation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.BotMessage{}) {
			return
		}

		m.responser.UpdateRecordResponse(c, m.store.AdminPanel().BotMessages(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	@Method GET
	@Path admin/messages
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/messages/messages_api.md#selection

	Получение лимитированного объема записей из таблицы `bot_messages`

	# TESTED
*/
func (m *ModMessage) GetMessagesSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c, m.store.AdminPanel().BotMessages(), &models.BotMessageSelection{})
}

/*
	@Method POST
	@Path admin/message
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/messages/messages_api.md#create

	Создать запись в таблице `bot_messages`

	# TESTED
*/
func (m *ModMessage) CreateNewMessageHandler(c *gin.Context) {
	r := &models.BotMessage{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.Validation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.BotMessage{}) {
			return
		}

		m.responser.CreateRecordResponse(c, m.store.AdminPanel().BotMessages(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}
