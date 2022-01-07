package bills

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
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
func (m *ModBills) GetAllBillsHandler(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidPathParams)
		return
	}

	b, err := m.store.AdminPanel().Bills().All(int64(chatID))
	m.responser.RecordResponse(c, b, err)
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
	chatID, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidPathParams)
		return
	}

	if obj := m.responser.RecordHandler(c, &models.Bill{ChatID: int64(chatID)}); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Bill{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.store.AdminPanel().Bills(), obj)
		return
	}
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
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
	}

	r := &models.Bill{ID: id}
	m.responser.RecordResponse(c, r, m.store.AdminPanel().Bills().FindById(r))
}
