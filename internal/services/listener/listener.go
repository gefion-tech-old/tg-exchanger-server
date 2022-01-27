package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/plugins"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
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
		_, err := l.snapshot(cfg)
		if err != nil {
			return err
		}

		// Получение истории транзакций
		if err := l.checker(); err != nil {
			return err
		}

		<-t.C
	}
}

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

	utils.SetSuccessStep(GetAllMerachants)

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
			utils.SetSuccessStep(fmt.Sprintf("%s %s", DecodeParams, m.Name))

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
				continue
			}

			utils.SetSuccessStep(fmt.Sprintf("Ping account %s", m.Name))
			c.Merchants.Whitebit = append(c.Merchants.Whitebit, p.(*models.WhitebitOptionParams))
		}

	}

	return c, nil
}

func (l *Listener) checker() error {

	fmt.Printf("Код выполнен - %v\n", time.Now().Format(time.UnixDate))

	return nil
}
