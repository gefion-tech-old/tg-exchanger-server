package guard

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type Guard struct {
	redis     *redisstore.AppRedisDictionaries
	secrets   *config.SecretsConfig
	responser utils.ResponserI
	logger    utils.LoggerI
}

type GuardI interface {
	IsAuth() gin.HandlerFunc
	IsAdmin() gin.HandlerFunc
	AuthTokenValidation() gin.HandlerFunc
	Logger(resource, action string) gin.HandlerFunc
}

func Init(r *redisstore.AppRedisDictionaries, s *config.SecretsConfig, resp utils.ResponserI, l utils.LoggerI) GuardI {
	return &Guard{
		redis:     r,
		secrets:   s,
		responser: resp,
		logger:    l,
	}
}
