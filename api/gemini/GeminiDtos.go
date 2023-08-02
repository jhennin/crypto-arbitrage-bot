package gemini

type payload struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
}

type PriceRecordGemini struct {
	Pair             string `json:"pair"`
	Price            string `json:"price"`
	PercentChange24h string `json:"percentChange24h"`
}

type errorGemini struct {
	result  string `json:"result"`
	reason  string `json:"reason"`
	message string `json:"message"`
}
