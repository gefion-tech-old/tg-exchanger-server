package ma

import (
	"encoding/json"
	"net/http"

	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
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

	if err := r.Validation(); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	var resp interface{}
	errs, _ := errgroup.WithContext(c)

	switch c.Param("service") {
	case AppType.MerchantAutoPayoutWhitebit:
		// Получение данных адреса
		errs.Go(func() error {
			b, err := m.pl.Whitebit.Merchant().CreateAdress(r)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(b.([]byte), &resp); err != nil {
				return err
			}

			return nil
		})

		// Создание заявки
		errs.Go(func() error {
			if err := m.repository.ExchangeRequest().Create(r); err != nil {
				return err
			}

			return nil
		})

	}

	if errs.Wait() != nil {
		m.responser.Error(c, http.StatusInternalServerError, errs.Wait())
		return
	}

	c.JSON(http.StatusOK, resp)
}
