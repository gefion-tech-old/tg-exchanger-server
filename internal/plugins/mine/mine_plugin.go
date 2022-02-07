package mine_plugin

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
)

type MinePlugin struct {
	merchant   interfaces.MerchantI
	autopayout interfaces.AutoPayoutI
}

func InitMinePlugin() interfaces.PluginI {
	return &MinePlugin{
		merchant:   InitMerchant(),
		autopayout: IniAutoPayout(),
	}
}

func (plugin *MinePlugin) Merchant() interfaces.MerchantI {
	if plugin.merchant != nil {
		return plugin.merchant
	}

	plugin.merchant = InitMerchant()
	return plugin.merchant
}

func (plugin *MinePlugin) AutoPayout() interfaces.AutoPayoutI {
	if plugin.autopayout != nil {
		return plugin.autopayout
	}

	plugin.autopayout = IniAutoPayout()
	return plugin.autopayout
}

func (plugin *MinePlugin) Ping(params interface{}) (interface{}, error) {
	return nil, nil
}

func (plugin *MinePlugin) History(params, body interface{}) (interface{}, error) {
	return nil, nil
}

func (plugin *MinePlugin) Balance(params, body interface{}) (interface{}, error) {
	return nil, nil
}

func (plugin *MinePlugin) GetOptionParams(options string) (interface{}, error) {
	return nil, nil
}
