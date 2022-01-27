package models

type ListenerState struct {
	Merchants ListenerMerchants
}

type ListenerMerchants struct {
	Whitebit []*WhitebitOptionParams
	Mine     []*MineOptionParams
}
