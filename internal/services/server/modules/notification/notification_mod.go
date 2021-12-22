package notification

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gin-gonic/gin"
)

type ModNotification struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cnf   *config.Config
}

type ModNotificationI interface {
	CreateNotificationHandler(c *gin.Context)
	GetAllNotificationsHandler(c *gin.Context)
	UpdateNotificationStatusHandler(c *gin.Context)
	DeleteNotificationHandler(c *gin.Context)
}

func InitModNotification(store db.SQLStoreI, redis *redisstore.AppRedisDictionaries, nsq nsqstore.NsqI, cnf *config.Config) ModNotificationI {
	return &ModNotification{
		store: store,
		redis: redis,
		nsq:   nsq,
		cnf:   cnf,
	}
}
