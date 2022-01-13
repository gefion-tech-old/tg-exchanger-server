package interfaces

type AppPlugins interface {
	Mine() PluginI
	Whitebit() PluginI
}

type PluginI interface {
	Merchant() MerchantI
	AutoPayout() AutoPayoutI
}

type MerchantI interface{}

type AutoPayoutI interface{}
