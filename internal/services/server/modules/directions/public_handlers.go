package directions

import "github.com/gin-gonic/gin"

func (m *ModDirections) GetDirectionHandler(c *gin.Context) {
	m.DirectionHandler(c)
}
