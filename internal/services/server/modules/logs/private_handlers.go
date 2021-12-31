package logs

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (m *ModLogs) DeleteLogRecordHandler(c *gin.Context) {
	m.responser.DeleteRecordResponse(c, m.repository, m.responser.DRRhelper(c, &models.LogRecord{}).(*models.LogRecord))
}

func (m *ModLogs) GetLogRecordsSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c, m.repository)
}

func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {}
