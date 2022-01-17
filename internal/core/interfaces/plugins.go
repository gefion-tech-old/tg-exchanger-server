package interfaces

type PluginI interface {
	Merchant() MerchantI
	AutoPayout() AutoPayoutI
}

type MerchantI interface {
	CreateAdress(d interface{}) (interface{}, error)
	GetHistory(d interface{}) (interface{}, error)
}

type AutoPayoutI interface{}
