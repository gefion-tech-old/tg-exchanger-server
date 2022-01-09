package guard

import (
	"context"
	"net/http"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
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
			g.responser.Error(c, http.StatusUnauthorized, err)
			return
		}

		_, err = g.redis.Auth.FetchAuth(tokenAuth)
		if err != nil {
			g.responser.Error(c, http.StatusUnauthorized, err)
			return
		}

		// Записываю в контекст структуру AccessDetails
		c.Request = c.Request.WithContext(context.WithValue(
			c.Request.Context(),
			CtxKeyToken,
			tokenAuth,
		))

		c.Next()
	}
}

func (g *Guard) IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := g.extractTokenMetadata(c.Request)
		if err != nil {
			g.responser.Error(c, http.StatusUnauthorized, err)
			return
		}

		if token.Role != AppType.AppRoleAdmin {
			go g.logger.NewRecord(&models.LogRecord{
				Service:  AppType.LogLevelAdmin,
				Module:   "GUARD",
				Info:     "Unauthorized access attempt",
				Username: &token.Username,
			})

			g.responser.Error(c, http.StatusForbidden, AppError.ErrNotEnoughRights)
			return
		}

		// Записываю в контекст структуру AccessDetails
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), CtxKeyToken, token))

		c.Next()
	}
}
