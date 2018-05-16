package service

import (
	"encoding/json"
	"fmt"

	"github.com/samoslab/priceservice/util"
)

/*[
    {
        "id": "skycoin",
        "name": "Skycoin",
        "symbol": "SKY",
        "rank": "87",
        "price_usd": "24.693",
        "price_btc": "0.0028395",
        "24h_volume_usd": "1666730.0",
        "market_cap_usd": "217640596.0",
        "available_supply": "8813858.0",
        "total_supply": "25000000.0",
        "max_supply": "100000000.0",
        "percent_change_1h": "-0.29",
        "percent_change_24h": "-2.87",
        "percent_change_7d": "14.97",
        "last_updated": "1526355251",
        "price_eur": "20.700364137",
        "24h_volume_eur": "1397234.75957",
        "market_cap_eur": "182450070.0"
    }
]*/
type (
	CoinMarketInfo struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		Symbol           string `json:"symbol"`
		Rank             string `json:"rank"`
		PriceUsd         string `json:"price_usd"`
		PriceBtc         string `json:"price_btc"`
		VolumeUsd24h     string `json:"24h_volume_usd"`
		MarketCapUsd     string `json:"market_cap_usd"`
		AvailableSupply  string `json:"available_supply"`
		TotalSupply      string `json:"total_supply"`
		MaxSupply        string `json:"max_supply"`
		PercentChange1h  string `json:"percent_change_1h"`
		PercentChange24h string `json:"percent_change_24h"`
		PercentChange7d  string `json:"percent_change_7d"`
		LastUpdated      string `json:"last_updated"`
		PriceCny         string `json:"price_cny"`
		VolumeCny24h     string `json:"24h_volume_cny"`
		MarketCapCny     string `json:"market_cap_cny"`
	}
)

func GetCoinPriceInfo(coinType string) (*CoinMarketInfo, error) {
	url := fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/%s/?convert=CNY", coinType)
	data, err := util.SendRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("get url err %v\n", err)
		return nil, err
	}
	priceInfo := []CoinMarketInfo{}
	if err := json.Unmarshal(data, &priceInfo); err != nil {
		fmt.Printf("unmarshal err %v\n", err)
		return nil, err
	}
	fmt.Printf("%s price btc %s, usd %s\n", coinType, priceInfo[0].PriceBtc, priceInfo[0].PriceUsd)
	return &priceInfo[0], nil
}
