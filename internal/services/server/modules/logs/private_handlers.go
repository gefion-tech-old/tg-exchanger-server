package logs

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

func (m *ModLogs) DeleteLogRecordHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	r := &models.LogRecord{ID: id}

	switch m.repository.Delete(r) {
	case nil:
		c.JSON(http.StatusOK, r)
	case sql.ErrNoRows:
		tools.ServErr(c, http.StatusNotFound, errors.ErrRecordNotFound)
		return
	default:
		tools.ServErr(c, http.StatusInternalServerError, err)
		return

	}
}

func (m *ModLogs) GetLogRecordsSelectionHandler(c *gin.Context) {}

func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {}
