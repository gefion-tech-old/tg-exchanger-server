package message

import (
	"net/http"
	"strconv"

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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	r := &models.BotMessage{ID: id}
	m.responser.RecordResponse(c, r, m.store.AdminPanel().BotMessages().Delete(r))
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	r.ID = id

	m.responser.Error(c, http.StatusUnprocessableEntity,
		r.UpdateBotMessageValidation(m.cnf.Users.Managers, m.cnf.Users.Developers),
	)

	m.responser.RecordResponse(c, r, m.store.AdminPanel().BotMessages().Update(r))
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

	m.responser.Error(c, http.StatusUnprocessableEntity,
		r.CreateBotMessageValidation(m.cnf.Users.Managers, m.cnf.Users.Developers),
	)

	m.responser.NewRecordResponse(c, r, m.store.AdminPanel().BotMessages().Create(r))
}
