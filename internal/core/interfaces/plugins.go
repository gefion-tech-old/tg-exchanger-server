package interfaces

type PluginI interface {
	Merchant() MerchantI
	AutoPayout() AutoPayoutI

	Ping(params interface{}) (interface{}, error)
	History(params, body interface{}) (interface{}, error)
	GetOptionParams(options string) (interface{}, error)
}

type MerchantI interface {
	CreateAdress(d, params interface{}) (interface{}, error)
}

type AutoPayoutI interface{}
