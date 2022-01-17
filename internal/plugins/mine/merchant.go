package mine_plugin

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
)

type MinePluginMerchant struct{}

func InitMerchant() interfaces.MerchantI {
	return &MinePluginMerchant{}
}

func (p *MinePluginMerchant) CreateAdress(d interface{}) (interface{}, error) {
	return nil, nil
}

func (p *MinePluginMerchant) GetHistory(d interface{}) (interface{}, error) {
	return nil, nil
}
