package listener

import (
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
)

func (l *Listener) snapshot(cfg *config.ListenerConfig) (*models.ListenerState, error) {
	c := &models.ListenerState{
		Merchants: models.ListenerMerchants{
			Whitebit: []*models.WhitebitOptionParams{},
			Mine:     []*models.MineOptionParams{},
		},
	}

	// Получение списка всех доступных мерчантов которые необходимо проверять
	arr, err := l.store.AdminPanel().MerchantAutopayout().GetAllMerchants()
	if err != nil {
		return nil, err
	}

	utils.SetSuccessStep(AppType.ListenerStepGetAllMerachants)

	for _, m := range arr {
		switch m.Service {
		case AppType.MerchantAutoPayoutWhitebit:
			// Декодирую опциональные параметры
			p, err := l.plugin.Whitebit.GetOptionParams(m.Options)
			if err != nil {
				return nil, fmt.Errorf("%s | ID: %d | %s",
					AppError.ErrFailedToDecodeParams.Error(),
					m.ID,
					err.Error(),
				)
			}
			utils.SetSuccessStep(AppType.SprintfStep("%s %s", DecodeParams, m.Name))

			// Пингую аккаунт
			b, err := l.plugin.Whitebit.Ping(p)
			if err != nil {
				return nil, fmt.Errorf("failed to ping account %s | %s", m.Name, err.Error())
			}

			var resp map[string]interface{}
			if err := json.Unmarshal([]byte(b.([]byte)), &resp); err != nil {
				return nil, fmt.Errorf("failed to ping account %s | %s", m.Name, err.Error())
			}

			if resp["code"] != nil {
				fmt.Printf("failed to ping account %s\n", m.Name)
				// TODO: Писать лог что не удалось установить соединение с этим аккаунтом
				continue
			}

			utils.SetSuccessStep(AppType.SprintfStep("Ping account %s", m.Name))
			c.Merchants.Whitebit = append(c.Merchants.Whitebit, p.(*models.WhitebitOptionParams))
		}

	}

	return c, nil
}
