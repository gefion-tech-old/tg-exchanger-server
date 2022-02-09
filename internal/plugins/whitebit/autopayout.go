package whitebit_plugin

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type WhitebitPluginAutoPayout struct{}

func IniAutoPayout() interfaces.AutoPayoutI {
	return &WhitebitPluginAutoPayout{}
}

func (p *WhitebitPluginAutoPayout) Payout(params, body interface{}) (interface{}, error) {
	b, err := SendRequest(
		params.(*models.WhitebitOptionParams),
		WhitebitbWithdrawPay,
		body.(map[string]interface{}),
	)
	if err != nil {
		return nil, err
	}
	return b, nil
}
