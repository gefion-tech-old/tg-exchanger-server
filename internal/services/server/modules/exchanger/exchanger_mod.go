package exchanger

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gin-gonic/gin"
)

type ModExchanger struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cnf   *config.Config
}

type ModExchangerI interface {
	CreateExchangerHandler(c *gin.Context)
	UpdateExchangerHandler(c *gin.Context)
	DeleteExchangerHandler(c *gin.Context)
	GetExchangerByNameHandler(c *gin.Context)
	GetExchangersSelectionHandler(c *gin.Context)
	GetExchangerDocumentHandler(c *gin.Context)
}

func InitModExchanger(store db.SQLStoreI, redis *redisstore.AppRedisDictionaries, nsq nsqstore.NsqI, cnf *config.Config) ModExchangerI {
	return &ModExchanger{
		store: store,
		redis: redis,
		nsq:   nsq,
		cnf:   cnf,
	}
}
