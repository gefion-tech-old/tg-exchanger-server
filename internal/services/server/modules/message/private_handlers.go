package message

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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
	m.responser.Record(c, r, m.store.AdminPanel().BotMessages().Delete(r))
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

	m.responser.Record(c, r, m.store.AdminPanel().BotMessages().Update(r))
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
	errs, _ := errgroup.WithContext(c)

	cArrM := make(chan []*models.BotMessage)
	cCount := make(chan *int)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Подсчет кол-ва сообщений в таблице
	errs.Go(func() error {
		defer close(cCount)
		c, err := m.store.AdminPanel().BotMessages().Count()
		if err != nil {
			return err
		}

		cCount <- &c
		return nil
	})

	// Достаю из БД запрашиваемые записи
	errs.Go(func() error {
		defer close(cArrM)
		arrM, err := m.store.AdminPanel().BotMessages().Selection(page, limit)
		if err != nil {
			return err
		}

		cArrM <- arrM
		return nil
	})

	arrM := <-cArrM
	count := <-cCount

	if arrM == nil || count == nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errs.Wait())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"current_page": page,
		"last_page":    math.Ceil(float64(*count) / float64(limit)),
		"data":         arrM,
	})
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

	m.responser.NewRecord(c, r, m.store.AdminPanel().BotMessages().Create(r))
}
