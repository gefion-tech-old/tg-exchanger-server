package directions

import (
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (m *ModDirections) CreateDirectionHandler(c *gin.Context) {
	m.DirectionHandler(c)
}

func (m *ModDirections) UpdateDirectionHandler(c *gin.Context) {
	m.DirectionHandler(c)
}

func (m *ModDirections) DeleteDirectionHandler(c *gin.Context) {
	m.DirectionHandler(c)
}

func (m *ModDirections) DirectionSelectionHandler(c *gin.Context) {
	s := &models.DirectionSelection{}
	if c.Query("status") != "" {
		b, err := strconv.ParseBool(c.Query("status"))
		if err != nil {
			m.responser.Error(c, http.StatusUnprocessableEntity, err)
			return
		}
		s = &models.DirectionSelection{Status: &b}
	}

	m.responser.SelectionResponse(c, m.repository.Directions(), s)
}
