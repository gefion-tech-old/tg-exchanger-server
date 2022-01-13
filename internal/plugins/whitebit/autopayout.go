package whitebit_plugin

import "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"

type WhitebitPluginAutoPayout struct{}

func IniAutoPayout() interfaces.AutoPayoutI {
	return &WhitebitPluginAutoPayout{}
}
