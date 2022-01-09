package logs

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type ModLogs struct {
	repository db.LoggerRepository
	cnf        *config.Config

	responser utils.ResponserI
}

type ModLogsI interface {
	CreateLogRecordHandler(c *gin.Context)
	DeleteLogRecordHandler(c *gin.Context)
	GetLogRecordsSelectionHandler(c *gin.Context)
	DeleteLogRecordsSelectionHandler(c *gin.Context)
}

func InitModLogs(r db.LoggerRepository, cnf *config.Config, responser utils.ResponserI) ModLogsI {
	return &ModLogs{
		repository: r,
		cnf:        cnf,
		responser:  responser,
	}
}
