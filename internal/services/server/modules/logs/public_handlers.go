package logs

import (
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path /log
	@Type PUBLIC
	@Documentation

	Создание записи в таблице `logs`
*/
func (m *ModLogs) CreateLogRecordHandler(c *gin.Context) {
	r := &models.LogRecord{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	switch r.Service {
	case static.L__SERVER:
		m.responser.Error(c, http.StatusUnprocessableEntity, r.InternalRecordValidation())

	case static.L__BOT:
		m.responser.Error(c, http.StatusUnprocessableEntity, r.InternalRecordValidation())

	case static.L__ADMIN:
		m.responser.Error(c, http.StatusUnprocessableEntity, r.AdminRecordValidation())

	default:
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	m.responser.NewRecordResponse(c, r, m.repository.Create(r))
}
