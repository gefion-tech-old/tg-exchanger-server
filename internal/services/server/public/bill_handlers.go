package public

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

/*
	@Method GET
	@Path /bot/user/:chat_id/bills
	@Type PUBLIC
	@Documentation

	Получить список всех имеющихся счетов у пользователя.

	# TESTED
*/
func (pr *PublicRoutes) getAllBillsHandler(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidPathParams)
		return
	}

	b, err := pr.store.User().Bills().All(int64(chatID))
	switch err {
	case nil:
		c.JSON(http.StatusOK, b)
	case sql.ErrNoRows:
		tools.ServErr(c, http.StatusNotFound, errors.ErrRecordNotFound)
		return
	default:
		tools.ServErr(c, http.StatusInternalServerError, err)
		return
	}
}

/*
	@Method DELETE
	@Path /bot/user/bill/:id
	@Type PUBLIC
	@Documentation

	В методе проверяется валидность переданых данных, если все ок
	и желаемый счет существует -> удаляю его из БД.

	# TESTED
*/
func (pr *PublicRoutes) deleteBillHandler(c *gin.Context) {
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

	_, err := pr.store.User().Bills().Delete(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{})

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": errors.ErrRecordNotFound.Error(),
		})
		return

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (pr *PublicRoutes) getBill(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	b, err := pr.store.User().Bills().FindById(&models.Bill{ID: uint(id)})
	switch err {
	case nil:
		c.JSON(http.StatusOK, b)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, errors.ErrRecordNotFound)
		return

	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}
}
