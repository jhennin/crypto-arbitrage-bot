package gemini

import (
	"crypto/hmac"
	"crypto/sha512"
	"cryptoArbitrageBot/bookkeeper"
	"cryptoArbitrageBot/internal"
	"cryptoArbitrageBot/internal/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const takerFeeGemini = .004

type GeminiClient struct {
	client *http.Client
}

/*
Create a new GeminiClient
*/
func NewClient() GeminiClient {
	return GeminiClient{
		client: internal.GetClient(),
	}
}

/*
Get crypto currency prices from Gemini. Specifically, get the price of BTCUSD, ETHUSD, LTCUSD, ETHBTC, LTCBTC, and LTCETH.
*/
func (c *GeminiClient) GetPrices() []bookkeeper.PriceRecord {

	priceFeedGemini, error := getPriceFeed(c.client)
	if error != nil {
		log.Fatal("ERROR getting signed prices")
	}

	var priceRecords = []bookkeeper.PriceRecord{}

	for _, aPrice := range priceFeedGemini {
		price, err := strconv.ParseFloat(aPrice.Price, 64)
		if err != nil {
			utils.Logger.Fatal(err.Error())
		}
		if aPrice.Pair == "BTCUSD" || aPrice.Pair == "ETHUSD" || aPrice.Pair == "LTCUSD" || aPrice.Pair == "ETHBTC" || aPrice.Pair == "LTCBTC" || aPrice.Pair == "LTCETH" {
			priceRecord := bookkeeper.PriceRecord{
				Uuid:                   uuid.New(),
				Currency:               aPrice.Pair,
				Price:                  price,
				Fee:                    takerFeeGemini,
				Exchange:               "Gemini",
				ArbitrageRecordUuid:    uuid.Nil,
				IsArbitrageOpportunity: false,
				Timestamp:              time.Now().Format(time.RFC3339),
			}
			priceRecords = append(priceRecords, priceRecord)
		}
	}

	utils.Logger.Debug("Retrieved Gemini prices...", zap.String("numberOfPriceRecords", strconv.Itoa(len(priceRecords))))
	bookkeeper.RecordPriceRecord(priceRecords...)
	return priceRecords
}

func getBtcPriceFromGeminiPriceFeed(priceFeed []PriceRecordGemini) string {
	for _, n := range priceFeed {
		if n.Pair == "BTCUSD" {
			return n.Price
		}
	}
	return ""
}

/*
Build a http.Request for Gemini.
*/
func requestBuilder(now string, path string, method string) *http.Request {
	geminiURL, ok := viper.Get("GEMINI.URL").(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	key, ok := viper.Get("GEMINI.TEST.KEY").(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	secret, ok := viper.Get("GEMINI.TEST.SECRET").(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	payload := payload{Request: path, Nonce: now}
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}
	signature := hmac.New(sha512.New384, []byte(secret))
	signature.Write(encodedPayload)
	xGeminiSignature := hex.EncodeToString(signature.Sum(nil))

	request, err := http.NewRequest(method, geminiURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Length", "0")
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("X-GEMINI-APIKEY", key)
	request.Header.Set("X-GEMINI-PAYLOAD", base64.StdEncoding.EncodeToString(encodedPayload))
	request.Header.Set("X-GEMINI-SIGNATURE", xGeminiSignature)
	request.Header.Set("Cache-Control", "no-cache")

	log.Printf("Built Gemini request.")
	return request
}

/*
Get the Gemini price feed.
*/
func getPriceFeed(client *http.Client) ([]PriceRecordGemini, error) {
	var priceRecordGemini []PriceRecordGemini
	errorGemini := errorGemini{}

	u, _ := url.ParseRequestURI(viper.Get("GEMINI.URL").(string))
	u.Path = "/v1/pricefeed"
	urlString := u.String()

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &priceRecordGemini)
	if err != nil {
		log.Println("Failed to unmarshal the following response body:\n\n", string(body)+"\nERROR UNMARSHALLING:", err)

		err = json.Unmarshal(body, &errorGemini)
		if err != nil {
			log.Println("Failed to unmarshal the following response body:\n\n", string(body)+"\nERROR UNMARSHALLING:", err)
		}
	}

	if errorGemini.message != "" {
		return priceRecordGemini, errors.New(errorGemini.message)
	}
	return priceRecordGemini, nil
}
