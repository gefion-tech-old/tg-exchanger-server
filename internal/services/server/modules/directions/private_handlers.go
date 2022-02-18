package directions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/* Внешние методя для направлений обмена */

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
		s.Status = &b
	}

	m.responser.SelectionResponse(c, m.repository.Directions(), s)
}

/* Внешние методя для мерчантов и автовыплат для конкретного направдения обмена */

func (m *ModDirections) CreateDirectionMaHandler(c *gin.Context) {
	m.DirectionMaHandler(c)
}

func (m *ModDirections) UpdateDirectionMaHandler(c *gin.Context) {
	m.DirectionMaHandler(c)
}

func (m *ModDirections) GetDirectionMaHandler(c *gin.Context) {
	m.DirectionMaHandler(c)
}

func (m *ModDirections) DeleteDirectionMaHandler(c *gin.Context) {
	m.DirectionMaHandler(c)
}

func (m *ModDirections) DirectionMaSelectionHandler(c *gin.Context) {
	s := &models.DirectionMASelection{}
	id, err := strconv.Atoi(c.Query("did"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}
	s.DirectionID = id
	fmt.Println(111)
	m.responser.SelectionResponse(c, m.repository.Directions().Ma(), s)
}
