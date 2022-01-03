package logs

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

/*
	@Method DELETE
	@Path /log
	@Type PRIVATE
	@Documentation

	Удаление записи из бд `logs`
*/
func (m *ModLogs) DeleteLogRecordHandler(c *gin.Context) {
	if obj := m.responser.RecordHandler(c, &models.LogRecord{}).(*models.LogRecord); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.LogRecord{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.repository, obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}

/*
	@Method GET
	@Path /logs
	@Type PRIVATE
	@Documentation

	Выборка по условиям из таблицы `logs`
*/
func (m *ModLogs) GetLogRecordsSelectionHandler(c *gin.Context) {
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

	if err := m.responser.DateHandler(c, c.Query("from"), c.Query("to")); err != nil {
		return
	}

	errs, _ := errgroup.WithContext(c)

	cArr := make(chan []*models.LogRecord)
	cCount := make(chan *int)

	errs.Go(func() error {
		defer close(cCount)

		c, err := m.repository.CountWithCustomFilters(c.Query("user"), c.Query("from"), c.Query("to"))
		if err != nil {
			return err
		}

		cCount <- &c
		return nil
	})

	errs.Go(func() error {
		defer close(cArr)

		arr, err := m.repository.SelectionWithCustomFilters(page, limit, c.Query("user"), c.Query("from"), c.Query("to"))
		if err != nil {
			return err
		}

		cArr <- arr
		return nil
	})

	arr := <-cArr
	count := <-cCount

	if arr == nil || count == nil {
		m.responser.Error(c, http.StatusInternalServerError, errs.Wait())
		return
	}

	m.responser.SelectionResponseObj(c, arr, page, limit, *count)
}

/*
	@Method DELETE
	@Path /logs
	@Type PRIVATE
	@Documentation

	Удаление записей по условиям из таблицы `logs`
*/
func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {
	if err := m.responser.DateHandler(c, c.Query("from"), c.Query("to")); err != nil {
		return
	}

	arr, err := m.repository.DeleteSelection(c.Query("from"), c.Query("to"))
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, arr)
}
