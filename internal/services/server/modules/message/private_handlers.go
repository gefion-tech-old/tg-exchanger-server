package message

import (
	"net/http"
	"reflect"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method DELETE
	@Path admin/message/:id
	@Type PRIVATE
	@Documentation

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

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}

/*
	@Method PUT
	@Path admin/message/:id
	@Type PRIVATE
	@Documentation

	Обновить запись в таблице `bot_messages`

	# TESTED
*/
func (m *ModMessage) UpdateBotMessageHandler(c *gin.Context) {
	r := &models.BotMessage{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r,
		r.UpdateBotMessageValidation(m.cnf.Users.Managers, m.cnf.Users.Developers),
	); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.BotMessage{}) {
			return
		}

		m.responser.UpdateRecordResponse(c, m.store.AdminPanel().BotMessages(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}

/*
	@Method GET
	@Path admin/messages
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `bot_messages`

	# TESTED
*/
func (m *ModMessage) GetMessagesSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c, m.store.AdminPanel().BotMessages())
}

/*
	@Method POST
	@Path admin/message
	@Type PRIVATE
	@Documentation

	Создать запись в таблице `bot_messages`

	# TESTED
*/
func (m *ModMessage) CreateNewMessageHandler(c *gin.Context) {
	r := &models.BotMessage{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r,
		r.CreateBotMessageValidation(m.cnf.Users.Managers, m.cnf.Users.Developers),
	); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.BotMessage{}) {
			return
		}

		m.responser.CreateRecordResponse(c, m.store.AdminPanel().BotMessages(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}
