package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/plugins"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"golang.org/x/sync/errgroup"
)

type Listener struct {
	store  db.SQLStoreI
	nsq    nsqstore.NsqI
	plugin *plugins.AppPlugins
	logger utils.LoggerI
}

type ListenerI interface {
	Listen(ctx context.Context, cfg *config.ListenerConfig) error
}

func InitListener(s db.SQLStoreI, q nsqstore.NsqI, p *plugins.AppPlugins, l utils.LoggerI) ListenerI {
	return &Listener{
		store:  s,
		nsq:    q,
		plugin: p,
		logger: l,
	}
}

func (l *Listener) Listen(ctx context.Context, cfg *config.ListenerConfig) error {
	for {
		t := time.NewTimer(time.Duration(cfg.Interval) * time.Second)

		// Получение актуального списка проверяемых аккаунтов
		state, err := l.snapshot(cfg)
		if err != nil {
			return err
		}

		// Массив историй транзакций со всех аккаунтов на Whitebit
		cWhitebitHistoryArr := make(chan []*models.WhitebitHistory)
		cExchangeRequestsArr := make(chan []*models.ExchangeRequest)

		var whitebitHistoryArr []*models.WhitebitHistory
		var exchangeRequestsArr []*models.ExchangeRequest

		// Сбор истории транзакции со всех акаунтов всех мерчантов
		{
			errs, _ := errgroup.WithContext(ctx)

			// Получение списка всех актуальных заявок
			errs.Go(func() error {
				defer close(cExchangeRequestsArr)

				arr, err := l.store.AdminPanel().ExchangeRequest().GetAllByStatus(AppType.ExchangeRequestNew)
				if err != nil {
					return err
				}

				cExchangeRequestsArr <- arr
				return nil
			})

			// Получаю истории транзакций всех whitebit аккаутов
			errs.Go(func() error {
				defer close(cWhitebitHistoryArr)
				arr := []*models.WhitebitHistory{}

				for _, merchant := range state.Merchants.Whitebit {
					history, err := l.checker(merchant)
					if err != nil {
						return err
					}

					arr = append(arr, history)
				}

				cWhitebitHistoryArr <- arr
				return nil
			})

			whitebitHistoryArr = <-cWhitebitHistoryArr
			exchangeRequestsArr = <-cExchangeRequestsArr

			if errs.Wait() != nil {
				fmt.Println(errs.Wait())
			}
		}

		// Анализ истории транзакций всех аккаунтов всех мерчантов
		{
			errs, _ := errgroup.WithContext(ctx)

			// Анализ истории всех транзакций со всех аккаунтов на whitebit
			errs.Go(func() error {

				// rHistory -> Запись из истории транзакеций
				// rRequest -> Запись в таблице заявок
				for _, account := range whitebitHistoryArr {
					for _, rHistory := range account.Records {
						for _, rRequest := range exchangeRequestsArr {
							if rHistory.Address == rRequest.Address {
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
								} else if rRequest.TransferredAmount < rRequest.ExpectedAmount {
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

								if err := l.store.AdminPanel().ExchangeRequest().Update(rRequest); err != nil {
									return err
								}

								utils.SetSuccessStep("Request processed")
								return nil
							}
						}
					}
				}

				return nil
			})

			if errs.Wait() != nil {
				fmt.Println(errs.Wait())
			}
		}

		<-t.C
	}
}

func (l *Listener) checker(m *models.WhitebitOptionParams) (*models.WhitebitHistory, error) {
	time.Sleep(time.Duration(1 * time.Second))

	b, err := l.plugin.Whitebit.History(m, AppType.BaseWhitebitGetHistoryBody)
	if err != nil {
		// TODO: Писать лог что не удалось установить соединение с этим аккаунтом
		return nil, err
	}

	var history models.WhitebitHistory
	if err := json.Unmarshal(b.([]byte), &history); err != nil {
		return nil, err
	}

	return &history, nil
}
