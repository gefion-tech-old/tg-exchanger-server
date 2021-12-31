package logs

import (
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (m *ModLogs) DeleteLogRecordHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	r := &models.LogRecord{ID: id}
	m.responser.RecordResponse(c, r, m.repository.Delete(r))
}

func (m *ModLogs) GetLogRecordsSelectionHandler(c *gin.Context) {}

func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {}
