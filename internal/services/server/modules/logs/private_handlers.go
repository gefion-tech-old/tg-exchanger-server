package logs

import (
	"net/http"
	"reflect"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

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

func (m *ModLogs) GetLogRecordsSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c, m.repository)
}

func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {}
