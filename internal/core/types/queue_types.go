package ctypes

// Типы событий отправляемых в очередь NSQ
var (
	QueueEventVerificationCode    = "verification_code"
	QueueEventSkipOperation       = "skip_operation"
	QueueEventConfirmationRequest = "confirmation_req"
	QueueEventExchangeError       = "exchange_error"
)
