package user

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type ModUsers struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cfg   *config.Config

	responser utils.ResponserI
	logger    utils.LoggerI
}

type ModUsersI interface {
	UserInBotRegistrationHandler(c *gin.Context)
	UserGenerateCodeHandler(c *gin.Context)
	UserInAdminRegistrationHandler(c *gin.Context)
	UserInAdminAuthHandler(c *gin.Context)
	UserRefreshToken(c *gin.Context)

	LogoutHandler(c *gin.Context)
}

func InitModUsers(
	store db.SQLStoreI,
	redis *redisstore.AppRedisDictionaries,
	nsq nsqstore.NsqI,
	cfg *config.Config,
	responser utils.ResponserI,
	l utils.LoggerI,
) ModUsersI {
	return &ModUsers{
		store: store,
		redis: redis,
		nsq:   nsq,
		cfg:   cfg,

		responser: responser,
		logger:    l,
	}
}
