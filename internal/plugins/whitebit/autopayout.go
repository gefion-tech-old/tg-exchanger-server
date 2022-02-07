package whitebit_plugin

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type WhitebitPluginAutoPayout struct{}

func IniAutoPayout() interfaces.AutoPayoutI {
	return &WhitebitPluginAutoPayout{}
}

func (p *WhitebitPluginAutoPayout) Payout(data *models.WhitebitWithdrawRequest) (interface{}, error) {
	return nil, nil
}
