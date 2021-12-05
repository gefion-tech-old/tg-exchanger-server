package guard

import (
	"context"
	"net/http"

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

		// if errFetch := g.redis.Get(tokenAuth.AccessUuid).Err(); errFetch != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": errors.ErrTokenInvalid,
		// 	})
		// 	c.Abort()
		// 	return
		// }

		// Записываю в контекст структуру AccessDetails
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), CtxKeyToken, tokenAuth))

		c.Next()
	}
}
