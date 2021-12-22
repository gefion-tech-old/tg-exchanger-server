package bills

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gin-gonic/gin"
)

type ModBills struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cnf   *config.Config
}

type ModBillsI interface {
	GetBillHandler(c *gin.Context)
	DeleteBillHandler(c *gin.Context)
	GetAllBillsHandler(c *gin.Context)

	RejectBillHandler(c *gin.Context)
	CreateBillHandler(c *gin.Context)
}

func InitModBills(store db.SQLStoreI, redis *redisstore.AppRedisDictionaries, nsq nsqstore.NsqI, cnf *config.Config) ModBillsI {
	return &ModBills{
		store: store,
		redis: redis,
		nsq:   nsq,
		cnf:   cnf,
	}
}
