package models

type WhitebitApiHistory struct {
	Account struct {
		Address string `json:"address"`
	} `json:"account"`

	Required struct {
		FixedFee string `json:"fixedFee"`
		FlexFee  struct {
			MaxFee  string `json:"maxFee"`
			MinFee  string `json:"minFee"`
			Percent string `json:"percent"`
		} `json:"flexFee"`

		MaxAmount string `json:"maxAmount"`
		MinAmount string `json:"minAmount"`
	} `json:"required"`
}

type WhitebitHistory struct {
	Limit   int                     `json:"limit"`
	Offset  int                     `json:"offset"`
	Records []WhitebitHistoryRecord `json:"records"`
	Total   int                     `json:"total"`
}

type WhitebitHistoryRecord struct {
	Address         string      `json:"address"`
	Amount          string      `json:"amount"`
	Centralized     bool        `json:"centralized"`
	CreatedAt       int         `json:"createdAt"`
	Currency        string      `json:"currency"`
	Description     interface{} `json:"description"`
	Fee             string      `json:"fee"`
	Memo            string      `json:"memo"`
	Method          int         `json:"method"`
	Network         interface{} `json:"network"`
	Status          int         `json:"status"`
	Ticker          string      `json:"ticker"`
	TransactionHash string      `json:"transactionHash"`
	UniqueId        interface{} `json:"uniqueId"`
}
