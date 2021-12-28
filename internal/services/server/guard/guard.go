package guard

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gin-gonic/gin"
)

type Guard struct {
	redis   *redisstore.AppRedisDictionaries
	secrets *config.SecretsConfig
}

type GuardI interface {
	IsAuth() gin.HandlerFunc
	IsAdmin() gin.HandlerFunc
	AuthTokenValidation() gin.HandlerFunc
}

func Init(r *redisstore.AppRedisDictionaries, s *config.SecretsConfig) GuardI {
	return &Guard{
		redis:   r,
		secrets: s,
	}
}
