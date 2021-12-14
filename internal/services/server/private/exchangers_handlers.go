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

func (pr *PrivateRoutes) createExchanger(c *gin.Context) {
	req := &models.Exchanger{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if err := req.ExchangerValidation(); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

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

func (pr *PrivateRoutes) updateExchanger(c *gin.Context) {
	req := &models.Exchanger{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if err := req.ExchangerValidationFull(); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

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

func (pr *PrivateRoutes) deleteExchanger(c *gin.Context) {
	req := &models.Exchanger{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if req.ID < 1 {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	e, err := pr.store.Manager().Exchanger().Delete(req)
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

func (pr *PrivateRoutes) getExchanger(c *gin.Context) {
	req := &models.Exchanger{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if req.ID < 1 {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	e, err := pr.store.Manager().Exchanger().Get(req)
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

func (pr *PrivateRoutes) getAllExchangers(c *gin.Context) {
	errs, _ := errgroup.WithContext(c)

	cArrE := make(chan []*models.Exchanger)
	cCount := make(chan *int)

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

	// Подсчет кол-ва уведомлений в таблице
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
