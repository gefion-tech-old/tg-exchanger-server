package plugins

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
)

type AppPlugins struct {
	Mine     interfaces.PluginI
	Whitebit interfaces.PluginI
}

func InitAppPlugins(m, wb interfaces.PluginI) *AppPlugins {
	return &AppPlugins{
		Mine:     m,
		Whitebit: wb,
	}
}
