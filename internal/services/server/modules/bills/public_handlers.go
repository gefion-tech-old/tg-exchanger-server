package bills

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
func (pr *ModBills) GetAllBillsHandler(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, errors.ErrInvalidPathParams)
		return
	}

	b, err := pr.store.AdminPanel().Bills().All(int64(chatID))
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
func (m *ModBills) DeleteBillHandler(c *gin.Context) {
	r := &models.Bill{}
	if err := c.ShouldBindJSON(r); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	if err := r.BillValidation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	m.responser.Record(c, r, m.store.AdminPanel().Bills().Delete(r))
}

/*
	@Method DELETE
	@Path /bot/user/bill/:id
	@Type PUBLIC
	@Documentation

	Получить запись из таблицы `bills`


	# TESTED
*/
func (m *ModBills) GetBillHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	r := &models.Bill{ID: id}
	m.responser.Record(c, r, m.store.AdminPanel().Bills().FindById(r))
}
