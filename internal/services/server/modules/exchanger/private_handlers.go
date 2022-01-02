package exchanger

import (
	"net/http"
	"reflect"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path admin/exchangers
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `exchangers`

	# TESTED
*/
func (m *ModExchanger) GetExchangersSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c, m.store.AdminPanel().Exchanger(), nil)
}

/*
	@Method DELETE
	@Path admin/exchanger/:id
	@Type PRIVATE
	@Documentation

	Удалить запись в таблице `exchangers`

	# TESTED
*/
func (m *ModExchanger) DeleteExchangerHandler(c *gin.Context) {
	if obj := m.responser.RecordHandler(c, &models.Exchanger{}); obj != nil {
		// Проверю, удалось ли записать структуру или была поймана ошибка
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Exchanger{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.store.AdminPanel().Exchanger(), obj)
	}

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}

/*
	@Method PUT
	@Path admin/exchanger/:id
	@Type PRIVATE
	@Documentation

	Обновить запись в таблице `exchangers`

	# TESTED
*/
func (m *ModExchanger) UpdateExchangerHandler(c *gin.Context) {
	// Декодирование
	r := &models.Exchanger{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.ExchangerUpdateValidation()); obj != nil {
		// Проверю, удалось ли записать структуру или была поймана ошибка
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Exchanger{}) {
			return
		}

		m.responser.UpdateRecordResponse(c, m.store.AdminPanel().Exchanger(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}

/*
	@Method POST
	@Path admin/exchanger
	@Type PRIVATE
	@Documentation

	Создать запись в таблице `exchangers`

	# TESTED
*/
func (m *ModExchanger) CreateExchangerHandler(c *gin.Context) {
	// Декодирование
	r := &models.Exchanger{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.ExchangerCreateValidation()); obj != nil {
		// Проверю, удалось ли записать структуру или была поймана ошибка
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Exchanger{}) {
			return
		}

		m.responser.CreateRecordResponse(c, m.store.AdminPanel().Exchanger(), obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}
