package whitebit_plugin_test

import (
	"fmt"
	"testing"

	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	whitebit_plugin "github.com/gefion-tech/tg-exchanger-server/internal/plugins/whitebit"
	"github.com/stretchr/testify/assert"
)

func Test_Plugin_Whitebit_Merchant_PrepareBodyForCreateAdress(t *testing.T) {
	testCases := []struct {
		data            *models.MerchantNewAdress
		expectedNetwork string
		expectedTicker  string
	}{
		{
			data: &models.MerchantNewAdress{
				Ticker: "USDTTRC20",
			},
			expectedNetwork: AppType.CurrencyNetworkTRC20,
			expectedTicker:  "USDT",
		},
		// {
		// 	data: &models.MerchantNewAdress{
		// 		Ticker: "BTC",
		// 	},
		// 	expectedNetwork: "",
		// 	expectedTicker:  "BTC",
		// },
		{
			data: &models.MerchantNewAdress{
				Ticker: "USDTOMNI",
			},
			expectedNetwork: AppType.CurrencyNetworkOMNI,
			expectedTicker:  "USDT",
		},
		{
			data: &models.MerchantNewAdress{
				Ticker: "USDTERC20",
			},
			expectedNetwork: AppType.CurrencyNetworkERC20,
			expectedTicker:  "USDT",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d\n", i), func(t *testing.T) {
			body := whitebit_plugin.PrepareBodyForCreateAdress(tc.data)

			assert.Equal(t, tc.expectedNetwork, body["network"])
			assert.Equal(t, tc.expectedTicker, body["ticker"])

		})
	}
}
