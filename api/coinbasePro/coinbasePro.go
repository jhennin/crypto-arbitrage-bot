package coinbasePro

import (
	"crypto/hmac"
	"crypto/sha256"
	"cryptoArbitrageBot/bookkeeper"
	"cryptoArbitrageBot/internal"
	"cryptoArbitrageBot/internal/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const takerFeeCoinbase = .006

type CoinbaseProClient struct {
	client *http.Client
}

type CoinbaseProError struct {
	Msg string
}

func (e *CoinbaseProError) Error() string {
	return fmt.Sprintf("Error using the Coinbase Pro API.")
}

func NewClient() CoinbaseProClient {
	return CoinbaseProClient{
		internal.GetClient(),
	}
}

func (c *CoinbaseProClient) GetPrices() ([]bookkeeper.PriceRecord, *CoinbaseProError) {

	aTicker, productTickerErr := getProductTicker(c.client, "BTC-USD") //TODO program is crashing here
	if productTickerErr != nil {
		utils.Logger.Error("Error getting BTC-USD price from Coinbase Pro.", zap.Error(productTickerErr))
	}
	btcPrice, err := strconv.ParseFloat(aTicker.Price, 64)
	if err != nil {
		utils.Logger.Error("Error converting BTC-USD price from Coinbase Pro.", zap.Error(err))

	}

	aTicker, productTickerErr = getProductTicker(c.client, "ETH-USD")
	if productTickerErr != nil {
		utils.Logger.Error("Error getting ETH-USD price from Coinbase Pro.", zap.Error(productTickerErr))
	}
	ethPrice, err := strconv.ParseFloat(aTicker.Price, 64)
	if err != nil {
		utils.Logger.Error("Error converting ETH-USD price from Coinbase Pro.", zap.Error(err))
	}

	aTicker, productTickerErr = getProductTicker(c.client, "LTC-USD")
	if productTickerErr != nil {
		utils.Logger.Error("Error getting LTC-USD price from Coinbase Pro.", zap.Error(productTickerErr))
	}
	ltcPrice, err := strconv.ParseFloat(aTicker.Price, 64)
	if err != nil {
		utils.Logger.Error("Error converting LTC-USD price from Coinbase Pro.", zap.Error(err))
	}

	aTicker, productTickerErr = getProductTicker(c.client, "ETH-BTC")
	if productTickerErr != nil {
		utils.Logger.Error("Error getting LTC-USD price from Coinbase Pro.", zap.Error(productTickerErr))
	}
	ethBtcPrice, err := strconv.ParseFloat(aTicker.Price, 64)
	if err != nil {
		utils.Logger.Error("Error converting ETH-BTC price from Coinbase Pro.", zap.Error(err))
	}

	aTicker, productTickerErr = getProductTicker(c.client, "LTC-BTC")
	if productTickerErr != nil {
		utils.Logger.Error("Error getting LTC-USD price from Coinbase Pro.", zap.Error(productTickerErr))
	}
	ltcBtcPrice, err := strconv.ParseFloat(aTicker.Price, 64)
	if err != nil {
		utils.Logger.Error("Error converting LTC-BTC price from Coinbase Pro.", zap.Error(err))
	}

	prices := []bookkeeper.PriceRecord{
		{
			Uuid:                   uuid.New(),
			Currency:               "BTCUSD",
			Price:                  btcPrice,
			Fee:                    takerFeeCoinbase,
			Exchange:               "Coinbase",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
		{
			Uuid:                   uuid.New(),
			Currency:               "ETHUSD",
			Price:                  ethPrice,
			Fee:                    takerFeeCoinbase,
			Exchange:               "Coinbase",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
		{
			Uuid:                   uuid.New(),
			Currency:               "LTCUSD",
			Price:                  ltcPrice,
			Fee:                    takerFeeCoinbase,
			Exchange:               "Coinbase",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
		{
			Uuid:                   uuid.New(),
			Currency:               "ETHBTC",
			Price:                  ethBtcPrice,
			Fee:                    takerFeeCoinbase,
			Exchange:               "Coinbase",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
		{
			Uuid:                   uuid.New(),
			Currency:               "LTCBTC",
			Price:                  ltcBtcPrice,
			Fee:                    takerFeeCoinbase,
			Exchange:               "Coinbase",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
	}
	utils.Logger.Debug("Retrieved coinbase prices...", zap.String("prices", strconv.Itoa(len(prices))))
	if len(prices) == 0 {
		return nil, &CoinbaseProError{Msg: "No prices returned from Coinbase Pro API."}
	}
	bookkeeper.RecordPriceRecord(prices...)
	return prices, nil
}

func requestBuilder(now string, method string, path string, body string) *http.Request {
	var key = viper.Get("COINBASE_PRO.TEST.KEY").(string)
	var secret = viper.Get("COINBASE_PRO.TEST.SECRET").(string)
	var passphrase = viper.Get("COINBASE_PRO.TEST.PASSPHRASE").(string)

	prehashString := now + method + path + body
	hmacKey, _ := base64.StdEncoding.DecodeString(secret)
	signature := hmac.New(sha256.New, hmacKey)
	signature.Write([]byte(prehashString))
	cbAccessSignature := base64.StdEncoding.EncodeToString(signature.Sum(nil))

	u, _ := url.ParseRequestURI(viper.Get("COINBASE_PRO.URL").(string))
	u.Path = path
	urlString := u.String()

	request, err := http.NewRequest(method, urlString, nil)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("CB-ACCESS-SIGN", cbAccessSignature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", now)
	request.Header.Add("CB-ACCESS-KEY", key)
	request.Header.Add("CB-ACCESS-PASSPHRASE", passphrase)

	utils.Logger.Debug(fmt.Sprintf("Finished building Coinbase Pro API request for %s%s.", u.Host, u.Path))
	return request
}

func getProductTicker(client *http.Client, productId string) (ProductTicker, *CoinbaseProError) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	path := fmt.Sprintf("/products/%s/ticker", productId)
	productTicker := ProductTicker{}
	errorCoinbasePro := errorCoinbasePro{}

	resp, err := client.Do(requestBuilder(now, "GET", path, ""))
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			utils.Logger.Error("Request timed out!")
		} else {
			utils.Logger.Error(err.Error())
		}
		return productTicker, &CoinbaseProError{Msg: err.Error()}
	}
	defer resp.Body.Close() //TODO program is crashing here

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		utils.Logger.Error(err.Error())
	}

	err = json.Unmarshal(body, &productTicker)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("ERROR UNMARSHALLING 'productTicker' --> %v\n\n Response body: \n\n%v", err, body))
	}

	err = json.Unmarshal(body, &errorCoinbasePro)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("ERROR UNMARSHALLING 'errorCoinbasePro'--> %v\n\n Response body: \n\n%v", err, body))
	}

	if errorCoinbasePro.Message != "" {
		return productTicker, &CoinbaseProError{
			Msg: errorCoinbasePro.Message,
		}
	}
	return productTicker, nil

}

func getSignedPrices(client *http.Client) (SignedPrices, *CoinbaseProError) {

	now := strconv.FormatInt(time.Now().Unix(), 10)
	path := "/oracle"
	signedPrices := SignedPrices{}
	errorCoinbasePro := errorCoinbasePro{}

	resp, err := client.Do(requestBuilder(now, "GET", path, ""))
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		utils.Logger.Error(err.Error())
	}
	bodyString := string(body)
	fmt.Println(bodyString)

	err = json.Unmarshal(body, &signedPrices)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("ERROR UNMARSHALLING 'signedPrices' --> %v\n\n Response body: \n\n%v", err, body))
	}

	err = json.Unmarshal(body, &errorCoinbasePro)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("ERROR UNMARSHALLING 'errorCoinbasePro'--> %v\n\n Response body: \n\n%v", err, body))
	}

	if errorCoinbasePro.Message != "" {
		return signedPrices, &CoinbaseProError{
			Msg: errorCoinbasePro.Message,
		}
	}
	return signedPrices, nil
}
