package guard

import (
	"bytes"
	"io/ioutil"
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

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				g.responser.Error(c, http.StatusInternalServerError, err)
				return
			}
		}

		// Возвращаю io.ReadCloser в исходное состояние.
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		go g.logger.NewRecord(&models.LogRecord{
			Service:  AppType.LogTypeAdmin,
			Module:   resource,
			Info:     string(bodyBytes),
			Username: &token.Username,
		})

		c.Next()
	}
}
