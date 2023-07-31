package coinbasePro

import "time"

type Account struct {
	id             string `json:"id"`
	currency       string `json:"currency"`
	balance        string `json:"balance"`
	available      string `json:"available"`
	hold           string `json:"hold"`
	profileId      string `json:"profile_id"`
	tradingEnabled string `json:"trading_enabled"`
}

type errorCoinbasePro struct {
	Message string `json:"message"`
}

type ProductTicker struct {
	TradeId int       `json:"trade_id"`
	Price   string    `json:"price"`
	Size    string    `json:"size"`
	Time    time.Time `json:"time"`
	Bid     string    `json:"bid"`
	Ask     string    `json:"ask"`
	Volume  string    `json:"volume"`
}

type SignedPrices struct {
	Timestamp  string   `json:"timestamp"`
	Messages   []string `json:"messages"`
	Signatures []string `json:"signatures"`
	Prices     struct {
		Btc string `json:"BTC"`
		Eth string `json:"ETH"`
		Ltc string `json:"LTC"`
	} `json:"prices"`
}
