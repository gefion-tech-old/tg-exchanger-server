package public

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gin-gonic/gin"
)

type PublicRoutes struct {
	store   db.SQLStoreI
	redis   db.RedisStoreI
	router  *gin.Engine
	secrets *config.SecretsConfig
}

type PublicRoutesI interface {
	ConfigurePublicRouter(router *gin.RouterGroup)
}

// Конструктор модуля публичных маршрутов
func Init(store db.SQLStoreI, redis db.RedisStoreI, router *gin.Engine, secrets *config.SecretsConfig) PublicRoutesI {
	return &PublicRoutes{
		store:   store,
		redis:   redis,
		router:  router,
		secrets: secrets,
	}
}

// Метод конфигуратор всех публичных маршрутов
func (pr *PublicRoutes) ConfigurePublicRouter(router *gin.RouterGroup) {
	router.POST("/test", pr.testHandler)
}

func (pr *PublicRoutes) testHandler(c *gin.Context) {}
