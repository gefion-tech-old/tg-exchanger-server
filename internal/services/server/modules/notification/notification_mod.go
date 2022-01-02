package notification

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type ModNotification struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cnf   *config.Config

	responser utils.ResponserI
	logger    utils.LoggerI
}

type ModNotificationI interface {
	CreateNotificationHandler(c *gin.Context)
	GetNotificationsSelectionHandler(c *gin.Context)
	UpdateNotificationStatusHandler(c *gin.Context)
	DeleteNotificationHandler(c *gin.Context)
	NewNotificationsCheckHandler(c *gin.Context)
}

func InitModNotification(store db.SQLStoreI, redis *redisstore.AppRedisDictionaries, nsq nsqstore.NsqI, cnf *config.Config, l utils.LoggerI, responser utils.ResponserI) ModNotificationI {
	return &ModNotification{
		store: store,
		redis: redis,
		nsq:   nsq,
		cnf:   cnf,

		responser: responser,
		logger:    l,
	}
}
