package whitebit_plugin

import (
	"strings"

	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type WhitebitPluginMerchant struct{}

func InitMerchant() interfaces.MerchantI {
	return &WhitebitPluginMerchant{}
}

// Создать адрес для принятия денег
func (p *WhitebitPluginMerchant) CreateAdress(d, params interface{}) (interface{}, error) {
	b, err := SendRequest(
		params.(*models.WhitebitOptionParams),
		WhitebitCreateNewAddress,
		PrepareBodyForCreateAdress(d.(*models.ExchangeRequest)),
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func PrepareBodyForCreateAdress(data *models.ExchangeRequest) map[string]interface{} {
	var network string
	networks := []string{AppType.CurrencyNetworkTRC20, AppType.CurrencyNetworkOMNI, AppType.CurrencyNetworkERC20}

	for i := 0; i < len(networks); i++ {
		if strings.Contains(data.ExchangeFrom, networks[i]) {
			network = networks[i]
		}
	}

	if network != "" {
		var t string
		if network == AppType.CurrencyNetworkOMNI {
			t = data.ExchangeFrom[:len(network)]
		} else {
			t = data.ExchangeFrom[:len(network)-1]
		}

		return map[string]interface{}{
			"ticker":  t,
			"network": network,
		}
	}

	return map[string]interface{}{
		"ticker": data.ExchangeFrom,
	}
}
