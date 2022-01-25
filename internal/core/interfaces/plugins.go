package interfaces

type PluginI interface {
	Merchant() MerchantI
	AutoPayout() AutoPayoutI

	GetOptionParams(options string) (interface{}, error)
}

type MerchantI interface {
	CreateAdress(d, apiHandler interface{}) (interface{}, error)
	GetHistory(d interface{}) (interface{}, error)
}

type AutoPayoutI interface{}
