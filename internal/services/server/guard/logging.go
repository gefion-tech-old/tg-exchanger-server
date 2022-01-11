package guard

import (
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

		// // var body map[string]interface{}
		// // if err := c.ShouldBindJSON(&body); err != nil {
		// // 	fmt.Println(err)
		// // 	return
		// // }

		// // jsonString, err := json.Marshal(body)
		// // if err != nil {
		// // 	fmt.Println(err)
		// // }
		// b := c.Request.Body
		// jsonData, err := ioutil.ReadAll(c.Request.Body)

		// fmt.Println(string(jsonData))

		go g.logger.NewRecord(&models.LogRecord{
			Service:  AppType.LogTypeAdmin,
			Module:   resource,
			Info:     c.Request.Body,
			Username: &token.Username,
		})

		c.Next()
	}
}
