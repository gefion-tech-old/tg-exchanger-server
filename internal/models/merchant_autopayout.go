package models

type MerchantNewAdress struct {
	Service string `json:"service"`
	Ticker  string `json:"ticker"`
}

type WhitebitGetHistory struct {
	TransactionMethod int           `json:"transactionMethod"`
	Ticker            string        `json:"ticker"`
	Address           string        `json:"address"`
	UniqueId          string        `json:"uniqueId"`
	Limit             int           `json:"limit"`
	Offset            int           `json:"offset"`
	Status            []interface{} `json:"status"`
}
