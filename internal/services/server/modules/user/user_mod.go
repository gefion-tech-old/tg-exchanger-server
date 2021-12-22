package user

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gin-gonic/gin"
)

type ModUsers struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cnf   *config.Config
}

type ModUsersI interface {
	UserInBotRegistrationHandler(c *gin.Context)
	UserGenerateCodeHandler(c *gin.Context)
	UserInAdminRegistrationHandler(c *gin.Context)
	UserInAdminAuthHandler(c *gin.Context)
	UserRefreshToken(c *gin.Context)

	LogoutHandler(c *gin.Context)
}

func InitModUsers(store db.SQLStoreI, redis *redisstore.AppRedisDictionaries, nsq nsqstore.NsqI, cnf *config.Config) ModUsersI {
	return &ModUsers{
		store: store,
		redis: redis,
		nsq:   nsq,
		cnf:   cnf,
	}
}
