package private

import (
	"database/sql"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

func (pr *PrivateRoutes) createBill(c *gin.Context) {
	req := &models.Bill{}
	if err := c.ShouldBindJSON(req); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if err := req.BillValidation(); err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	bill, err := pr.store.User().Bills().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, bill)
	case sql.ErrNoRows:
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrAlreadyExists)
		return
	default:
		tools.ServErr(c, http.StatusInternalServerError, err)
		return
	}
}
