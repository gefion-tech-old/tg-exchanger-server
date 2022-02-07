package ctypes

// Типы доступных сетей криптовалют
var (
	CurrencyNetworkOMNI  = "OMNI"
	CurrencyNetworkTRC20 = "TRC20"
	CurrencyNetworkERC20 = "ERC20"
)

// Типы доступных сервисов для мерчантов/автовыплат
var (
	MerchantAutoPayoutMine     = "mine"
	MerchantAutoPayoutWhitebit = "whitebit"
)

// Типы доступных состояний
var (
	UseAsMerchant   = 1
	UseAsAutoPayout = 2
)

var BaseWhitebitGetHistoryBody = map[string]interface{}{
	"limit":  100,
	"offset": 0,
}
