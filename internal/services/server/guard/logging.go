package guard

import (
	"fmt"
	"net/http"

	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
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
			Service:  AppType.LogTypeAdmin,
			Module:   resource,
			Info:     fmt.Sprintf("%s %s", action, resource),
			Username: &token.Username,
		})

		c.Next()
	}
}
