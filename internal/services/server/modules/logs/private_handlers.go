package logs

import (
	"net/http"
	"reflect"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
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
	m.responser.SelectionResponse(c, m.repository, &models.LogRecordSelection{
		Username: c.Query("user"),
		Service:  []string{c.Query("service")},
		DateFrom: c.Query("from"),
		DateTo:   c.Query("to"),
	})
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
