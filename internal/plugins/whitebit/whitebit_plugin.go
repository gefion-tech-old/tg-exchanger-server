package whitebit_plugin

import (
	"encoding/hex"
	"encoding/json"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	AppMath "github.com/gefion-tech/tg-exchanger-server/internal/core/math"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

type WhitebitPlugin struct {
	merchant   interfaces.MerchantI
	autopayout interfaces.AutoPayoutI

	cfg *config.PluginsConfig
}

func InitWhitebitPlugin(cfg *config.PluginsConfig) interfaces.PluginI {
	return &WhitebitPlugin{
		merchant:   InitMerchant(),
		autopayout: IniAutoPayout(),

		cfg: cfg,
	}
}

func (plugin *WhitebitPlugin) Merchant() interfaces.MerchantI {
	if plugin.merchant != nil {
		return plugin.merchant
	}

	plugin.merchant = InitMerchant()
	return plugin.merchant
}

func (plugin *WhitebitPlugin) AutoPayout() interfaces.AutoPayoutI {
	if plugin.autopayout != nil {
		return plugin.autopayout
	}

	plugin.autopayout = IniAutoPayout()
	return plugin.autopayout
}

func (plugin *WhitebitPlugin) GetOptionParams(options string) (interface{}, error) {
	var p models.WhitebitOptionParams
	dOptions, err := AppMath.AesDecrypt(options, hex.EncodeToString([]byte(plugin.cfg.AesKey)))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(dOptions), &p); err != nil {
		return nil, err
	}

	return &p, nil
}
