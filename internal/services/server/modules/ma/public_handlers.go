package ma

import (
	"database/sql"
	"encoding/json"
	"net/http"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/sync/errgroup"
)

func (m *ModMerchantAutoPayout) CreateNewAdressHandler(c *gin.Context) {
	r := &models.ExchangeRequest{
		Status: AppType.ExchangeRequestNew,
	}

	if err := c.ShouldBindJSON(&r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	{
		errs, _ := errgroup.WithContext(c)

		errs.Go(func() error {
			return r.Validation()
		})

		errs.Go(func() error {
			return validation.Validate(
				c.Param("service"),
				validation.Required,
				validation.In(
					AppType.MerchantAutoPayoutWhitebit,
					AppType.MerchantAutoPayoutMine,
				),
			)
		})

		if errs.Wait() != nil {
			m.responser.Error(c, http.StatusInternalServerError, errs.Wait())
			return
		}
	}

	// Поиск доступного аккаунта
	ma, err := m.repository.MerchantAutopayout().GetFistIfActive(c.Param("service"))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			m.responser.Error(c, http.StatusNotFound, AppError.ErrNoMerchantAutopatout)
			return
		default:
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}
	}

	// Прасинг опциональных параметров
	p, err := m.pl.Whitebit.GetOptionParams(ma.Options)
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	var resp map[string]interface{}

	switch c.Param("service") {
	case AppType.MerchantAutoPayoutWhitebit:
		// Получение данных адреса
		b, err := m.pl.Whitebit.Merchant().CreateAdress(r, p.(*models.WhitebitOptionParams))
		if err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}

		if err := json.Unmarshal(b.([]byte), &resp); err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}

		r.Address = resp["account"].(map[string]interface{})["address"].(string)

		// Создание заявки
		if err := m.repository.ExchangeRequest().Create(r); err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}
