package interfaces

import "github.com/gefion-tech/tg-exchanger-server/internal/config"

type AppPlugins interface {
	Mine() PluginI
	Whitebit(cfg *config.WhitebitConfig) PluginI
}

type PluginI interface {
	Merchant() MerchantI
	AutoPayout() AutoPayoutI
}

type MerchantI interface{}

type AutoPayoutI interface{}
