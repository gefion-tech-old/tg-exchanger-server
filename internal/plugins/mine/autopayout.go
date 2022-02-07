package mine_plugin

import "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"

type MinePluginAutoPayout struct{}

func IniAutoPayout() interfaces.AutoPayoutI {
	return &MinePluginAutoPayout{}
}

func (p *MinePluginAutoPayout) Payout(params, body interface{}) (interface{}, error) {
	return nil, nil
}
