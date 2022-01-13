package mine_plugin

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
)

type MinePluginMerchant struct{}

func InitMerchant() interfaces.MerchantI {
	return &MinePluginMerchant{}
}
