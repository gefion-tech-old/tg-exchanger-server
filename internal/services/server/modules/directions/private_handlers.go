package directions

import "github.com/gin-gonic/gin"

func (m *ModDirections) CreateDirectionHandler(c *gin.Context) {
	m.DirectionHandler(c)
}
