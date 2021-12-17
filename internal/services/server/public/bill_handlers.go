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

/*
	@Method POST
	@Path /bot/user/bill
	@Type PUBLIC
	@Documentation

	В методе проверяется валидность переданых данных, если все ок создается
	банковский счет закрепленный за конкретным пользователем.

	# TESTED
*/
func (pr *PublicRoutes) newBillHandler(c *gin.Context) {
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
