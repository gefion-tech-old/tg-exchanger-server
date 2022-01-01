package guard

import (
	"fmt"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (g *Guard) Logger(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := g.extractTokenMetadata(c.Request)
		if err != nil {
			g.responser.Error(c, http.StatusUnauthorized, err)
			return
		}

		go g.logger.NewRecord(&models.LogRecord{
			Service:  static.L__ADMIN,
			Module:   resource,
			Info:     fmt.Sprintf("%s %s", action, resource),
			Username: &token.Username,
		})

		c.Next()
	}
}
