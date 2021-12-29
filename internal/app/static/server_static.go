package static

// Роли
const (
	S__ROLE__ADMIN   = 2
	S__ROLE__MANAGER = 1
	S__ROLE__USER    = 0
)

// Типы логов
const (
	L__SERVER = 100
	L__BOT    = 200
	L__ADMIN  = 300
)

// Нобор статусов уведомлений
const (
	NTF__S__NEW        = 1
	NTF__S__IN_PROCESS = 2
	NTF__S__COMPLETED  = 3
)

// Нобор типов уведомлений
const (
	NTF__T__VERIFICATION   = 854
	NTF__T__EXCHANGE_ERROR = 855

	NTF__T__REQ_SUPPORT = 900
)
