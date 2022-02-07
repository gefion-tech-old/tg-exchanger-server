package listener

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"golang.org/x/sync/errgroup"
)

func (listener *Listener) snapshot(ctx context.Context, cfg *config.ListenerConfig) (*models.ListenerState, error) {
	c := &models.ListenerState{
		Merchants: &models.ListeningAccounts{
			Whitebit: []*models.WhitebitOptionParams{},
			Mine:     []*models.MineOptionParams{},
		},
	}
	errs, _ := errgroup.WithContext(ctx)

	// Получение списка всех доступных мерчантов
	errs.Go(func() error {
		arr, err := listener.store.AdminPanel().MerchantAutopayout().GetAllByServiceType(AppType.UseAsMerchant, true)
		if err != nil {
			return err
		}

		utils.SetSuccessStep(AppType.ListenerStepGetAllMerachants)

		for _, m := range arr {
			switch m.Service {
			case AppType.MerchantAutoPayoutWhitebit:
				p, err := listener.ping(m)
				if err != nil {
					return err
				}

				utils.SetSuccessStep(AppType.SprintfStep("Ping account %s", m.Name))
				c.Merchants.Whitebit = append(c.Merchants.Whitebit, p.(*models.WhitebitOptionParams))
			}
		}

		return nil
	})

	// Получение списка всех доступных автовыплат
	errs.Go(func() error {
		arr, err := listener.store.AdminPanel().MerchantAutopayout().GetAllByServiceType(AppType.UseAsAutoPayout, true)
		if err != nil {
			return err
		}

		utils.SetSuccessStep(AppType.ListenerStepGetAllAutopayouts)

		for _, m := range arr {
			switch m.Service {
			case AppType.MerchantAutoPayoutWhitebit:
				p, err := listener.ping(m)
				if err != nil {
					return err
				}

				utils.SetSuccessStep(AppType.SprintfStep("Ping account %s", m.Name))
				c.Merchants.Whitebit = append(c.Autopayouts.Whitebit, p.(*models.WhitebitOptionParams))
			}
		}

		return nil
	})

	return c, errs.Wait()
}

func (listener *Listener) ping(m *models.MerchantAutopayout) (interface{}, error) {
	// Декодирую опциональные параметры
	p, err := listener.plugin.Whitebit.GetOptionParams(m.Options)
	if err != nil {
		return nil, fmt.Errorf("%s | ID: %d | %s",
			AppError.ErrFailedToDecodeParams.Error(),
			m.ID,
			err.Error(),
		)
	}

	utils.SetSuccessStep(AppType.SprintfStep("%s %s", DecodeParams, m.Name))

	// Пингую аккаунт
	b, err := listener.plugin.Whitebit.Ping(p)
	if err != nil {
		return nil, fmt.Errorf("failed to ping account %s | %s", m.Name, err.Error())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(b.([]byte)), &resp); err != nil {
		return nil, fmt.Errorf("failed to ping account %s | %s", m.Name, err.Error())
	}

	if resp["code"] != nil {
		return nil, fmt.Errorf("failed to ping account %s", m.Name)
	}

	return p, nil
}
