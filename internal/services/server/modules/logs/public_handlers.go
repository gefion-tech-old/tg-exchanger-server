package logs

import (
	"fmt"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

func (m *ModLogs) CreateLogRecordHandler(c *gin.Context) {
	r := &models.LogRecord{}
	if err := c.ShouldBindJSON(r); err != nil {
		fmt.Println(err)
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	switch r.Service {
	case static.L__BOT:
		if err := r.InternalRecordValidation(); err != nil {
			tools.ServErr(c, http.StatusUnprocessableEntity, err)
			return
		}

	case static.L__ADMIN:
		if err := r.AdminRecordValidation(); err != nil {
			tools.ServErr(c, http.StatusUnprocessableEntity, err)
			return
		}

	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	m.responser.NewRecord(c, r, m.repository.Create(r))
}
