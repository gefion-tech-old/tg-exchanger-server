package private

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/guard"
	"github.com/gin-gonic/gin"
)

type PrivateRoutes struct {
	store   db.SQLStoreI
	redis   *redisstore.AppRedisDictionaries
	router  *gin.Engine
	secrets *config.SecretsConfig
	users   *config.UsersConfig
}

type PrivateRoutesI interface {
	ConfigurePrivateRouter(router *gin.RouterGroup, g guard.GuardI)
}

// Конструктор модуля приватных маршрутов
func Init(
	store db.SQLStoreI,
	redis *redisstore.AppRedisDictionaries,
	router *gin.Engine,
	secrets *config.SecretsConfig,
	users *config.UsersConfig) PrivateRoutesI {
	return &PrivateRoutes{
		store:   store,
		redis:   redis,
		router:  router,
		secrets: secrets,
		users:   users,
	}
}

// Метод конфигуратор всех публичных маршрутов
func (pr *PrivateRoutes) ConfigurePrivateRouter(router *gin.RouterGroup, g guard.GuardI) {
	admin := router.Group("/admin")
	admin.POST("/logout", g.AuthTokenValidation(), g.IsAuth(), pr.logoutHandler)

	/* Работа с конкретными ресурсами */

	admin.POST("/message", g.AuthTokenValidation(), g.IsAuth(), pr.createNewBotMessageHandler)
	admin.GET("/message", pr.getBotMessageHandler)
	admin.PUT("/message", g.AuthTokenValidation(), g.IsAuth(), pr.updateAllBotMessageHandler)
	admin.DELETE("/message", g.AuthTokenValidation(), g.IsAuth(), pr.deleteBotMessageHandler)

	/* Работа с общим списокм конкретного ресурсами */

	admin.GET("/messages", pr.getAllBotMessageHandler)
}
