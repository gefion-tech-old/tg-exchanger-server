package whitebit_plugin

import (
	"strings"

	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type WhitebitPluginMerchant struct {
	provider *apiHelper
}

func InitMerchant(p *apiHelper) interfaces.MerchantI {
	return &WhitebitPluginMerchant{
		provider: p,
	}
}

// Создать адрес для принятия денег
func (p *WhitebitPluginMerchant) CreateAdress(d interface{}) (interface{}, error) {
	b, err := p.provider.SendRequest(
		WhitebitCreateNewAddress,
		PrepareBodyForCreateAdress(d.(*models.MerchantNewAdress)),
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (p *WhitebitPluginMerchant) GetHistory(d interface{}) (interface{}, error) {
	data := d.(*models.WhitebitGetHistory)
	b, err := p.provider.SendRequest(WhitebitHistory, map[string]interface{}{
		"transactionMethod": data.TransactionMethod,
		"limit":             data.Limit,
		"offset":            data.Offset,
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}

func PrepareBodyForCreateAdress(data *models.MerchantNewAdress) map[string]interface{} {
	var network string
	networks := []string{AppType.CurrencyNetworkTRC20, AppType.CurrencyNetworkOMNI, AppType.CurrencyNetworkERC20}

	for i := 0; i < len(networks); i++ {
		if strings.Contains(data.Ticker, networks[i]) {
			network = networks[i]
		}
	}

	if network == "" {
		return map[string]interface{}{
			"ticker": "BTC",
		}
	} else {
		var t string
		if network == AppType.CurrencyNetworkOMNI {
			t = data.Ticker[:len(network)]
		} else {
			t = data.Ticker[:len(network)-1]
		}

		return map[string]interface{}{
			"ticker":  t,
			"network": network,
		}
	}
}
