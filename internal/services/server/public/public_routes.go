package public

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gin-gonic/gin"
)

type PublicRoutes struct {
	store   db.SQLStoreI
	redis   *redisstore.AppRedisDictionaries
	router  *gin.Engine
	secrets *config.SecretsConfig
	users   *config.UsersConfig
}

type PublicRoutesI interface {
	ConfigurePublicRouter(router *gin.RouterGroup)
}

// Конструктор модуля публичных маршрутов
func Init(store db.SQLStoreI, redis *redisstore.AppRedisDictionaries, router *gin.Engine, secrets *config.SecretsConfig, users *config.UsersConfig) PublicRoutesI {
	return &PublicRoutes{
		store:   store,
		redis:   redis,
		router:  router,
		secrets: secrets,
		users:   users,
	}
}

// Метод конфигуратор всех публичных маршрутов
func (pr *PublicRoutes) ConfigurePublicRouter(router *gin.RouterGroup) {
	bot := router.Group("/bot")
	bot.POST("/registration", pr.userInBotRegistrationHandler)

	admin := router.Group("/admin")
	admin.POST("/registration/code", pr.userGenerateCodeHandler)
	admin.POST("/registration", pr.userInAdminRegistrationHandler)
}
