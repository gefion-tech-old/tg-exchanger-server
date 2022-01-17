package whitebit_plugin

var (
	WhitebitBalanceRoute       = "/api/v4/trade-account/balance"           // Получить баланс аккаунта
	WhitebitMainAccountAddress = "/api/v4/main-account/address"            // Получить адрес депозита криптовалюты
	WhitebitCreateNewAddress   = "/api/v4/main-account/create-new-address" // Сгенерировать адрес
	WhitebitHistory            = "/api/v4/main-account/history"            // Получить историю депозитов/снятий
)
