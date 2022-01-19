package logs

import (
	"net/http"
	"reflect"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method DELETE
	@Path /log
	@Type PRIVATE
	@Documentation

	Удаление записи из бд `logs`

	# TESTED
*/
func (m *ModLogs) DeleteLogRecordHandler(c *gin.Context) {
	if obj := m.responser.RecordHandler(c, &models.LogRecord{}); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.LogRecord{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.repository, obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	@Method GET
	@Path /logs
	@Type PRIVATE
	@Documentation

	Выборка по условиям из таблицы `logs`

	# TESTED
*/
func (m *ModLogs) GetLogRecordsSelectionHandler(c *gin.Context) {
	s := &models.LogRecordSelection{
		Username: c.Query("user"),
		Service:  []string{c.Query("service")},
		DateFrom: c.Query("from"),
		DateTo:   c.Query("to"),
	}

	m.responser.SelectionResponse(c, m.repository, s)
}

/*
	@Method DELETE
	@Path /logs
	@Type PRIVATE
	@Documentation

	Удаление записей по условиям из таблицы `logs`

	# TESTED
*/
func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {
	lrs := &models.LogRecordSelection{
		DateFrom: c.Query("from"),
		DateTo:   c.Query("to"),
	}

	if err := lrs.Validation(); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	arr, err := m.repository.DeleteSelection(lrs)
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, arr)
}
