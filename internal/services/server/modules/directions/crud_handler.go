package directions

import (
	"net/http"
	"reflect"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (m *ModDirections) DirectionHandler(c *gin.Context) {
	var r models.Direction
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&r); err != nil {
			m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
			return
		}

		if err := r.Validation(); err != nil {
			m.responser.Error(c, http.StatusUnprocessableEntity, err)
			return
		}
	}

	if obj := m.responser.RecordHandler(c, &r); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Direction{}) {
			return
		}

		switch c.Request.Method {
		case http.MethodPost:
			m.responser.CreateRecordResponse(c, m.repository.Directions(), obj)
			return
		case http.MethodGet:
			m.responser.GetRecordResponse(c, m.repository.Directions(), obj)
			return
		case http.MethodPut:
			m.responser.UpdateRecordResponse(c, m.repository.Directions(), obj)
			return
		case http.MethodDelete:
			m.responser.DeleteRecordResponse(c, m.repository.Directions(), obj)
			return
		}
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

func (m *ModDirections) DirectionMaHandler(c *gin.Context) {
	var r models.DirectionMA
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&r); err != nil {
			m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
			return
		}

		if err := r.Validation(); err != nil {
			m.responser.Error(c, http.StatusUnprocessableEntity, err)
			return
		}
	}

	if obj := m.responser.RecordHandler(c, &r); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.DirectionMA{}) {
			return
		}

		switch c.Request.Method {
		case http.MethodPost:
			m.responser.CreateRecordResponse(c, m.repository.Directions().Ma(), obj)
			return
		case http.MethodGet:
			m.responser.GetRecordResponse(c, m.repository.Directions().Ma(), obj)
			return
		case http.MethodPut:
			m.responser.UpdateRecordResponse(c, m.repository.Directions().Ma(), obj)
			return
		case http.MethodDelete:
			m.responser.DeleteRecordResponse(c, m.repository.Directions().Ma(), obj)
			return
		}
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}
