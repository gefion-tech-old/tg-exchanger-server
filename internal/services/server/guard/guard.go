package guard

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

type Guard struct {
	redis   *redis.Client
	secrets *config.SecretsConfig
}

type GuardI interface {
	IsAuth() gin.HandlerFunc
	AuthTokenValidation() gin.HandlerFunc
}

func Init(r *redis.Client, s *config.SecretsConfig) GuardI {
	return &Guard{
		redis:   r,
		secrets: s,
	}
}
