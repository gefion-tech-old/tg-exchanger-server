package logs

import (
	"net/http"

	CoreErrors "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path /log
	@Type PUBLIC
	@Documentation

	Создание записи в таблице `logs`

	# TESTED
*/
func (m *ModLogs) CreateLogRecordHandler(c *gin.Context) {
	r := &models.LogRecord{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, CoreErrors.ErrInvalidBody)
		return
	}

	if err := r.Validation(); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
	}

	m.responser.NewRecordResponse(c, r, m.repository.Create(r))
}
