package whitebit_plugin

import "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"

type WhitebitPluginMerchant struct{}

func InitMerchant() interfaces.MerchantI {
	return &WhitebitPluginMerchant{}
}
