package logs

import (
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
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

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
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
	if err := m.responser.DateHandler(c, c.Query("from"), c.Query("to")); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Валидация номера сервиса
	if c.Query("service") != "" {
		if err := validation.Validate(c.Query("service"),
			validation.In(
				strconv.Itoa(static.L__BOT),
				strconv.Itoa(static.L__SERVER),
				strconv.Itoa(static.L__ADMIN),
			),
		); err != nil {
			m.responser.Error(c, http.StatusUnprocessableEntity, err)
		}
	}

	// Валидация имени имени пользователя
	if c.Query("user") != "" {
		if err := validation.Validate(c.Query("user"),
			validation.Match(regexp.MustCompile(static.REGEX__NAME)),
		); err != nil {
			m.responser.Error(c, http.StatusUnprocessableEntity, err)
		}
	}

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

	# TESTED
*/
func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {
	if err := m.responser.DateHandler(c, c.Query("from"), c.Query("to")); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	arr, err := m.repository.DeleteSelection(c.Query("from"), c.Query("to"))
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, arr)
}
