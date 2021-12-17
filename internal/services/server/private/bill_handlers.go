package private

import (
	"database/sql"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (pr *PrivateRoutes) createBill(c *gin.Context) {
	req := &models.Bill{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	if err := req.BillValidation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	bill, err := pr.store.User().Bills().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, bill)
	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": errors.ErrAlreadyExists.Error(),
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}
