package private

import (
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/guard"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path admin/logout
	@Type PRIVATE

	При валидных данных токена в Redis удаляется
	токен используемый в текущей сессии.


*/
func (pr *PrivateRoutes) logout(c *gin.Context) {
	// Извлекаю метаданные JWT
	ctxToken := c.Request.Context().Value(guard.CtxKeyToken).(*models.AccessDetails)

	// Удаляю токен
	_, err := pr.redis.Auth.Del(ctxToken.AccessUuid).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrTokenInvalid,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully logged out",
	})
}
