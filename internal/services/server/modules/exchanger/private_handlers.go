package exchanger

import (
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path admin/exchangers
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `exchangers`

	# TESTED
*/
func (m *ModExchanger) GetExchangersSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c, m.store.AdminPanel().Exchanger())
}

/*
	@Method DELETE
	@Path admin/exchanger/:id
	@Type PRIVATE
	@Documentation

	Удалить запись в таблице `exchangers`

	# TESTED
*/
func (m *ModExchanger) DeleteExchangerHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	r := &models.Exchanger{ID: id}
	m.responser.RecordResponse(c, r, m.store.AdminPanel().Exchanger().Delete(r))
}

/*
	@Method PUT
	@Path admin/exchanger/:id
	@Type PRIVATE
	@Documentation

	Обновить запись в таблице `exchangers`

	# TESTED
*/
func (m *ModExchanger) UpdateExchangerHandler(c *gin.Context) {
	// Декодирование
	r := &models.Exchanger{}
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

	// Валидация
	m.responser.Error(c, http.StatusUnprocessableEntity, r.ExchangerUpdateValidation())

	// Операция с БД
	m.responser.RecordResponse(c, r, m.store.AdminPanel().Exchanger().Update(r))
}

/*
	@Method POST
	@Path admin/exchanger
	@Type PRIVATE
	@Documentation

	Создать запись в таблице `exchangers`

	# TESTED
*/
func (m *ModExchanger) CreateExchangerHandler(c *gin.Context) {
	// Декодирование
	r := &models.Exchanger{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидация
	m.responser.Error(c, http.StatusUnprocessableEntity, r.ExchangerCreateValidation())

	// Операция записи в БД
	m.responser.NewRecordResponse(c, r, m.store.AdminPanel().Exchanger().Create(r))
}
