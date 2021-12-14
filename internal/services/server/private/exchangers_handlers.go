package private

import (
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

func (pr *PrivateRoutes) createExchanger(c *gin.Context) {
	req := &models.Exchanger{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

}

func (pr *PrivateRoutes) updateExchanger(c *gin.Context) {}

func (pr *PrivateRoutes) deleteExchanger(c *gin.Context) {}

func (pr *PrivateRoutes) getExchanger(c *gin.Context) {}

func (pr *PrivateRoutes) getAllExchangers(c *gin.Context) {}
