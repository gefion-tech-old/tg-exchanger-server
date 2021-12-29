package exchanger

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
	@Method POST
	@Path admin/exchangers
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `exchangers`

	# TESTED
*/
func (m *ModExchanger) GetExchangersSelectionHandler(c *gin.Context) {
	errs, _ := errgroup.WithContext(c)

	cArrE := make(chan []*models.Exchanger, 1)
	cCount := make(chan *int, 1)

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

	// Достаю из БД запрашиваемые записи
	errs.Go(func() error {
		defer close(cArrE)
		arrE, err := m.store.AdminPanel().Exchanger().Selection(page, limit)
		if err != nil {
			return err
		}

		cArrE <- arrE
		return nil
	})

	// Подсчет кол-ва записей в таблице
	errs.Go(func() error {
		defer close(cCount)
		c, err := m.store.AdminPanel().Exchanger().Count()
		if err != nil {
			return err
		}

		cCount <- &c
		return nil
	})

	arrE := <-cArrE
	count := <-cCount

	if arrE == nil || count == nil {
		m.responser.Error(c, http.StatusInternalServerError, errs.Wait())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"current_page": page,
		"last_page":    math.Ceil(float64(*count) / float64(limit)),
		"data":         arrE,
	})
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
	m.responser.Record(c, r, m.store.AdminPanel().Exchanger().Delete(r))
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
	m.responser.Record(c, r, m.store.AdminPanel().Exchanger().Update(r))
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
	m.responser.NewRecord(c, r, m.store.AdminPanel().Exchanger().Create(r))
}
