package mocks

import AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"

var MerchantAutopayout = map[string]interface{}{
	"name":         "main_accout",
	"service":      AppType.MerchantAutoPayoutWhitebit,
	"service_type": AppType.UseAsMerchant,
	"options":      "{}",
	"status":       true,
	"message_id":   1,
	"created_by":   "I0HuKc",
}
