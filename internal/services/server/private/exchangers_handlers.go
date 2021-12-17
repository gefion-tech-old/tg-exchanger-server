package private

import (
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

/*
	@Method POST
	@Path admin/exchanger
	@Type PRIVATE
	@Documentation

	Создать запись в таблице `exchangers`
*/
func (pr *PrivateRoutes) createExchanger(c *gin.Context) {
	// Декодирование
	req := &models.Exchanger{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидация
	if err := req.ExchangerValidation(); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Операция записи в БД
	e, err := pr.store.Manager().Exchanger().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, e)
		return
	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}
}

/*
	@Method PUT
	@Path admin/exchanger/:id
	@Type PRIVATE
	@Documentation

	Обновить запись в таблице `exchangers`
*/
func (pr *PrivateRoutes) updateExchanger(c *gin.Context) {
	// Декодирование
	req := &models.Exchanger{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидация
	if err := req.ExchangerValidation(); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Операция обновления записи в БД
	e, err := pr.store.Manager().Exchanger().Update(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, e)
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, errors.ErrRecordNotFound)
		return

	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}
}

/*
	@Method DELETE
	@Path admin/exchanger/:id
	@Type PRIVATE
	@Documentation

	Удалить запись в таблице `exchangers`
*/
func (pr *PrivateRoutes) deleteExchanger(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Операция удаления записи из БД
	e, err := pr.store.Manager().Exchanger().Delete(&models.Exchanger{ID: id})
	switch err {
	case nil:
		c.JSON(http.StatusOK, e)
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, errors.ErrRecordNotFound)
		return

	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}
}

/*
	@Method POST
	@Path admin/exchangers/:name
	@Type PRIVATE
	@Documentation

	Получение одной записи из таблицы `exchangers`
*/
func (pr *PrivateRoutes) getExchangerByName(c *gin.Context) {
	// Операция получения записи из БД
	e, err := pr.store.Manager().Exchanger().GetByName(&models.Exchanger{Name: c.Param("name")})
	switch err {
	case nil:
		c.JSON(http.StatusOK, e)
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, errors.ErrRecordNotFound)
		return

	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}
}

/*
	@Method POST
	@Path admin/exchangers
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `exchangers`
*/
func (pr *PrivateRoutes) getAllExchangers(c *gin.Context) {
	errs, _ := errgroup.WithContext(c)

	cArrE := make(chan []*models.Exchanger, 1)
	cCount := make(chan *int, 1)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Достаю из БД запрашиваемые записи
	errs.Go(func() error {
		defer close(cArrE)
		arrE, err := pr.store.Manager().Exchanger().GetSlice(page * limit)
		if err != nil {
			return err
		}

		cArrE <- arrE
		return nil
	})

	// Подсчет кол-ва записей в таблице
	errs.Go(func() error {
		defer close(cCount)
		c, err := pr.store.Manager().Exchanger().Count()
		if err != nil {
			return err
		}

		cCount <- &c
		return nil
	})

	arrE := <-cArrE
	count := <-cCount

	if arrE == nil || count == nil {
		tools.ServErr(c, http.StatusInternalServerError, errs.Wait())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"current_page": page,
		"last_page":    math.Ceil(float64(*count) / float64(limit)),
		"data":         arrE[(page-1)*limit : tools.UpperThreshold(page, limit, *count)],
	})
}