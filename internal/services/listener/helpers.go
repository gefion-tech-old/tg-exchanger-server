package listener

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

func (l *Listener) moneyHasBeenSent(rHistory models.WhitebitHistoryRecord, u models.UserFromBotRequest) error {
	text := fmt.Sprintf("✅Чек операции `%s`✅\n\nСумма перевода: *%s %s*\nАдрес: *%s*\nХеш: `%s`",
		rHistory.UniqueId,
		rHistory.Amount,
		rHistory.Ticker,
		rHistory.Address,
		rHistory.TransactionHash,
	)
	// Подготовка уведомления
	payload, err := json.Marshal(prepareNotification(u, text))
	if err != nil {
		return err
	}

	// Отправка уведомления
	return l.nsq.Publish(AppType.TopicBotMessages, payload)
}

func (l *Listener) amountLessThanExpected(u models.UserFromBotRequest) error {
	// Подготовка уведомления
	payload, err := json.Marshal(prepareNotification(
		u,
		"❗️Отмена операции обмена❗️\n\nПереведенная сумма меньше ожидаемой.",
	))
	if err != nil {
		return err
	}

	// Отправка уведомления
	return l.nsq.Publish(AppType.TopicBotMessages, payload)
}

func (l *Listener) amountMoreThanExpected(u models.UserFromBotRequest) error {
	// Подготовка уведомления
	payload, err := json.Marshal(prepareNotification(
		u,
		"❗️Отмена операции обмена❗️\n\nПереведенная сумма больше ожидаемой. Не переживайте, наши менеджеры скоро свяжутся с вами.",
	))
	if err != nil {
		return err
	}

	// Отправка уведомления
	return l.nsq.Publish(AppType.TopicBotMessages, payload)
}

func prepareNotification(u models.UserFromBotRequest, text string) map[string]interface{} {
	return map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  u.ChatID,
			"username": u.Username,
		},
		"message": map[string]interface{}{
			"type": AppType.QueueEventExchangeError,
			"text": text,
		},
		"created_at": time.Now().UTC().Format(core.DateStandart),
	}
}
