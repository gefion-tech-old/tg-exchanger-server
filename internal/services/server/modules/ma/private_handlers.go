package ma

import (
	"net/http"
	"reflect"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path admin/merchant-autopayout
	@Type PRIVATE
	@Documentation

	Создать запись в таблице `merchant_autopayout`

	# TESTED
*/
func (m *ModMerchantAutoPayout) CreateMerchantAutopayoutHandler(c *gin.Context) {
	r := &models.MerchantAutopayout{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.Validation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.MerchantAutopayout{}) {
			return
		}

		m.responser.CreateRecordResponse(c, m.repository, obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	@Method PUT
	@Path admin/merchant-autopayout/:id
	@Type PRIVATE
	@Documentation

	Обновить запись в таблице `merchant_autopayout`

	# TESTED
*/
func (m *ModMerchantAutoPayout) UpdateMerchantAutopayoutHandler(c *gin.Context) {
	r := &models.MerchantAutopayout{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.Validation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.MerchantAutopayout{}) {
			return
		}

		m.responser.UpdateRecordResponse(c, m.repository, obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	@Method DELETE
	@Path admin/merchant-autopayout/:id
	@Type PRIVATE
	@Documentation

	Удалить запись в таблице `merchant_autopayout`

	# TESTED
*/
func (m *ModMerchantAutoPayout) DeleteMerchantAutopayoutHandler(c *gin.Context) {
	if obj := m.responser.RecordHandler(c, &models.MerchantAutopayout{}); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.MerchantAutopayout{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.repository, obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	@Method DELETE
	@Path admin/merchant-autopayout/:id
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `merchant_autopayout`

	# TESTED
*/
func (m *ModMerchantAutoPayout) GetMerchantAutopayoutSelectionHandler(c *gin.Context) {
	s := &models.MerchantAutopayoutSelection{
		Service: []string{c.Query("service")},
	}

	m.responser.SelectionResponse(c, m.repository, s)
}
