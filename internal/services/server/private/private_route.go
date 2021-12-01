package private

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gin-gonic/gin"
)

type PrivateRoutes struct {
	store   db.SQLStoreI
	redis   db.RedisStoreI
	router  *gin.Engine
	secrets *config.SecretsConfig
}

type PrivateRoutesI interface {
	ConfigurePrivateRouter(router *gin.RouterGroup)
}

// Конструктор модуля приватных маршрутов
func Init(store db.SQLStoreI, redis db.RedisStoreI, router *gin.Engine, secrets *config.SecretsConfig) PrivateRoutesI {
	return &PrivateRoutes{
		store:   store,
		redis:   redis,
		router:  router,
		secrets: secrets,
	}
}

// Метод конфигуратор всех публичных маршрутов
func (pr *PrivateRoutes) ConfigurePrivateRouter(router *gin.RouterGroup) {
}
