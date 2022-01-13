package whitebit_plugin

import "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"

type WhitebitPlugin struct {
	merchant   interfaces.MerchantI
	autopayout interfaces.AutoPayoutI
}

func InitWhitebitPlugin() interfaces.PluginI {
	return &WhitebitPlugin{
		merchant:   InitMerchant(),
		autopayout: IniAutoPayout(),
	}
}

func (plugin *WhitebitPlugin) Merchant() interfaces.MerchantI {
	if plugin.merchant != nil {
		return plugin.merchant
	}

	plugin.merchant = InitMerchant()
	return plugin.merchant
}

func (plugin *WhitebitPlugin) AutoPayout() interfaces.AutoPayoutI {
	if plugin.autopayout != nil {
		return plugin.autopayout
	}

	plugin.autopayout = IniAutoPayout()
	return plugin.autopayout
}
