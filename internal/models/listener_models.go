package models

type ListenerState struct {
	Merchants   *ListeningAccounts
	Autopayouts *ListeningAccounts
}

type ListeningAccounts struct {
	Whitebit []*WhitebitOptionParams
	Mine     []*MineOptionParams
}
