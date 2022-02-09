package whitebit_plugin

var (
	WhitebitMainAccountAddress = "/api/v4/main-account/address"            // Получить адрес депозита криптовалюты
	WhitebitCreateNewAddress   = "/api/v4/main-account/create-new-address" // Сгенерировать адрес
	WhitebitHistory            = "/api/v4/main-account/history"            // Получить историю депозитов/снятий
	WhitebitbBalance           = "/api/v4/main-account/balance"            // Баланс текущего счета
	WhitebitbWithdraw          = "/api/v4/main-account/withdraw"           // Вывод средств
	WhitebitbWithdrawPay       = "/api/v4/main-account/withdraw-pay"       // Создать запрос на вывод с определенной суммой вывода (комиссия не включена)
)
