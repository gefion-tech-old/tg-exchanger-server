package ma

import (
	"encoding/json"
	"net/http"

	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (m *ModMerchantAutoPayout) CreateNewAdressHandler(c *gin.Context) {
	r := &models.MerchantNewAdress{
		Ticker:  c.Query("ticker"),
		Service: c.Param("service"),
	}

	switch r.Service {
	case AppType.MerchantAutoPayoutWhitebit:
		b, err := m.pl.Whitebit.Merchant().CreateAdress(r)
		if err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}

		var resp interface{}
		if err := json.Unmarshal(b.([]byte), &resp); err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, resp)
		return
	}
}
