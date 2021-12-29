package logs

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gin-gonic/gin"
)

type ModLogs struct {
	repository db.LoggerRepository
	cnf        *config.Config
}

type ModLogsI interface {
	CreateLogRecordHandler(c *gin.Context)
}

func InitModLogs(r db.LoggerRepository, cnf *config.Config) ModLogsI {
	return &ModLogs{
		repository: r,
		cnf:        cnf,
	}
}
