package message

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type ModMessage struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cfg   *config.Config

	responser utils.ResponserI
	logger    utils.LoggerI
}

type ModMessageI interface {
	DeleteBotMessageHandler(c *gin.Context)
	UpdateBotMessageHandler(c *gin.Context)
	GetMessagesSelectionHandler(c *gin.Context)
	CreateNewMessageHandler(c *gin.Context)

	GetMessageHandler(c *gin.Context)
}

func InitModMessage(
	store db.SQLStoreI,
	redis *redisstore.AppRedisDictionaries,
	nsq nsqstore.NsqI,
	cfg *config.Config,
	responser utils.ResponserI,
	l utils.LoggerI,
) ModMessageI {
	return &ModMessage{
		store: store,
		redis: redis,
		nsq:   nsq,
		cfg:   cfg,

		responser: responser,
		logger:    l,
	}
}
