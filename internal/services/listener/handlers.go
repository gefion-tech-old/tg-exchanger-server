package listener

import (
	"fmt"
	"strconv"

	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
)

/* Обработка событий биржи Whitebit */

// Метод обработки события вывода средств
func (l *Listener) handleWhitebitWithdrawAction(rHistory models.WhitebitHistoryRecord, rRequest *models.ExchangeRequest) error {
	fmt.Println(rHistory.Amount)
	return nil
}

// Метод обработки события нового депозита
func (l *Listener) handleWhitebitDepositAction(rHistory models.WhitebitHistoryRecord, rRequest *models.ExchangeRequest) error {
	if rHistory.Address == rRequest.Address {
		// Проверяю статус операции
		if rHistory.Status == 3 || rHistory.Status == 7 {
			// Получаю переведенную сумму
			amount, err := strconv.ParseFloat(rHistory.Amount, 64)
			if err != nil {
				return err
			}

			// Сохраняю сумму полученную от пользователя
			rRequest.TransferredAmount = amount
			rRequest.TransactionHash = &rHistory.TransactionHash

			if rRequest.TransferredAmount == rRequest.ExpectedAmount {
				// Если полученная сумма совпадает с ожидаемой суммой
				rRequest.Status = AppType.ExchangeRequestPaid
			} else if rRequest.TransferredAmount > rRequest.ExpectedAmount {
				// Если полученная сумма меньше ожидаемой суммы
				rRequest.Status = AppType.ExchangeRequestInvalidAmount

				// Отправка уведомления
				if err := l.amountLessThanExpected(rRequest.CreatedBy); err != nil {
					return err
				}
			} else {
				// Если полученная сумма больше ожидаемой суммы
				rRequest.Status = AppType.ExchangeRequestPaid

				// Отправка уведомления
				if err := l.amountMoreThanExpected(rRequest.CreatedBy); err != nil {
					return err
				}
			}

			// Если все ок, обновляю запись в БД
			if err := l.store.AdminPanel().ExchangeRequest().Update(rRequest); err != nil {
				return err
			}

			utils.SetSuccessStep("Request processed")
			return nil
		}
	}

	return nil
}
