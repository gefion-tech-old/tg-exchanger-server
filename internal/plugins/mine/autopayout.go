package mine_plugin

import "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"

type MinePluginAutoPayout struct{}

func IniAutoPayout() interfaces.AutoPayoutI {
	return &MinePluginAutoPayout{}
}
