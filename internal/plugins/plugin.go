package plugins

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	mp "github.com/gefion-tech/tg-exchanger-server/internal/plugins/mine"
	wp "github.com/gefion-tech/tg-exchanger-server/internal/plugins/whitebit"
)

type AppPlugins struct {
	mine     interfaces.PluginI
	whitebit interfaces.PluginI
}

func InitAppPlugins() interfaces.AppPlugins {
	return &AppPlugins{}
}

func (p *AppPlugins) Mine() interfaces.PluginI {
	if p.mine != nil {
		return p.mine
	}

	p.mine = mp.InitMinePlugin()
	return p.mine
}

func (p *AppPlugins) Whitebit() interfaces.PluginI {
	if p.whitebit != nil {
		return p.whitebit
	}

	p.whitebit = wp.InitWhitebitPlugin()
	return p.whitebit
}
