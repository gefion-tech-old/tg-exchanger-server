package ctypes

type ExchangeRequestStatus int

// Возможные статусы состояния заявки
var (
	// Новая заявка
	ExchangeRequestNew ExchangeRequestStatus = 100

	// Отмененная пользователем
	ExchangeRequestCanceled ExchangeRequestStatus = 200

	// Удаленная/Старая заявка
	// Старые заявки в этот статус кидаем,
	// на которые клиенты забили
	ExchangeRequestDeleted ExchangeRequestStatus = 300

	// Оплаченная заявка
	// Когда мерчант увидел, что клиент успешно перечислил
	// средства и они зашли на нашу платежную систему
	ExchangeRequestPaid ExchangeRequestStatus = 400

	// Заявка на проверке
	// Когда мерчант увидел, что клиент перечислил средства,
	// но не в полной сумме или не в том коде валют и т.д.
	// В этом случае наш менеджер должен будет вручную проверить,
	// что случилось с заявкой и вынести какое-то решение по ней
	ExchangeRequestChecked ExchangeRequestStatus = 500

	// Ошибка автовыплаты
	// Когда почему-то не сработала автовыплата.
	// Например, сайт платежной системы лег, или выдалась
	// ошибка по апи и т.д
	ExchangeRequestAutopayoutError ExchangeRequestStatus = 600

	// Выполненная заявка
	ExchangeRequestDone ExchangeRequestStatus = 700
)
