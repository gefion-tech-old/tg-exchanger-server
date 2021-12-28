package guard

import (
	"context"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

/*
	Финальный метод верификации токена доступа.
	Из полученного токена извлекается полезная нагрузка
	и проверяется метаданные JWT в redis
	Если все ок, далее в контексте передается структура
	AccessDetails для дальнейшей работы
*/
func (g *Guard) IsAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenAuth, err := g.extractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		_, err = g.redis.Auth.FetchAuth(tokenAuth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// Записываю в контекст структуру AccessDetails
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), CtxKeyToken, tokenAuth))

		c.Next()
	}
}

func (g *Guard) IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := g.extractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if token.Role != static.S__ROLE__ADMIN {
			tools.ServErr(c, http.StatusForbidden, errors.ErrNotEnoughRights)
		}

		// Записываю в контекст структуру AccessDetails
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), CtxKeyToken, token))

		c.Next()
	}
}
