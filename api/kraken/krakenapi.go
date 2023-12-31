package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"cryptoArbitrageBot/bookkeeper"
	"cryptoArbitrageBot/internal/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// APIURL is the official Kraken API Endpoint
var APIURL string

const (
	// APIVersion is the official Kraken API Version Number
	APIVersion = "0"
	// APIUserAgent identifies this library with the Kraken API
	APIUserAgent = "Kraken GO API Agent (https://github.com/beldur/kraken-go-api-client)"
)

func init() {
	APIURL = "https://api.kraken.com"
}

// List of valid public methods
var publicMethods = []string{
	"Assets",
	"AssetPairs",
	"Time",
	"Trades",
}

// List of valid private methods
var privateMethods = []string{
	"AddExport",
	"AddOrder",
	"Balance",
	"CancelOrder",
	"ClosedOrders",
	"DepositAddresses",
	"DepositMethods",
	"DepositStatus",
	"ExportStatus",
	"GetWebSocketsToken",
	"Ledgers",
	"OpenOrders",
	"OpenPositions",
	"QueryLedgers",
	"QueryOrders",
	"QueryTrades",
	"RemoveExport",
	"RetrieveExport",
	"TradeBalance",
	"TradesHistory",
	"TradeVolume",
	"WalletTransfer",
	"Withdraw",
	"WithdrawCancel",
	"WithdrawInfo",
	"WithdrawStatus",
}

// These represent the minimum order sizes for the respective coins
// Should be monitored through here: https://support.kraken.com/hc/en-us/articles/205893708-What-is-the-minimum-order-size-
const (
	takerFeeKraken = .0026
)

// KrakenApi represents a Kraken API Client connection
type KrakenApi = KrakenAPI

// KrakenAPI represents a Kraken API Client connection
type KrakenAPI struct {
	key    string
	secret string
	client *http.Client
}

// New creates a new Kraken API client
func New(key, secret string) *KrakenAPI {
	krakenAPI := KrakenAPI{
		key:    key,
		secret: secret,
		client: http.DefaultClient,
	}
	return &krakenAPI
}

// NewWithClient creates a new Kraken API client with custom http client
func NewWithClient(key, secret string, httpClient *http.Client) *KrakenAPI {
	kraken := New(key, secret)
	return kraken.WithClient(httpClient)
}

// WithClient adds an HTTP client into the KrakenAPI
func (api *KrakenAPI) WithClient(httpClient *http.Client) *KrakenAPI {
	api.client = httpClient
	return api
}

// Time returns the server's time
func (api *KrakenAPI) Time() (*TimeResponse, error) {
	resp, err := api.queryPublic("Time", nil, &TimeResponse{})
	if err != nil {
		return nil, err
	}

	return resp.(*TimeResponse), nil
}

// Assets returns the servers available assets
func (api *KrakenAPI) Assets() (*AssetsResponse, error) {
	resp, err := api.queryPublic("Assets", nil, &AssetsResponse{})
	if err != nil {
		return nil, err
	}

	return resp.(*AssetsResponse), nil
}

/*
	Returns the ticker for given comma separated pairs
*/
func (api *KrakenAPI) Ticker(pairs ...string) (*TickerResponse, error) {
	resp, err := api.queryPublic("Ticker", url.Values{
		"pair": {strings.Join(pairs, ",")},
	}, &TickerResponse{})
	if err != nil {
		return nil, err
	}

	return resp.(*TickerResponse), nil
}

/*
Get prices from Kraken. Specifically,
*/
func (api *KrakenAPI) GetPrices() []bookkeeper.PriceRecord {

	resp, err := api.Ticker()
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	btcAskPrice, err := strconv.ParseFloat(resp.XBTUSDT.Ask[0], 64)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	ethAskPrice, err := strconv.ParseFloat(resp.XETHZUSD.Ask[0], 64)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	ltcAskPrice, err := strconv.ParseFloat(resp.XLTCZUSD.Ask[0], 64)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	priceRecords := []bookkeeper.PriceRecord{
		{
			Uuid:                   uuid.New(),
			Currency:               "BTCUSD",
			Price:                  btcAskPrice,
			Fee:                    takerFeeKraken,
			Exchange:               "Kraken",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
		{
			Uuid:                   uuid.New(),
			Currency:               "ETHUSD",
			Price:                  ethAskPrice,
			Fee:                    takerFeeKraken,
			Exchange:               "Kraken",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
		{
			Uuid:                   uuid.New(),
			Currency:               "LTCUSD",
			Price:                  ltcAskPrice,
			Fee:                    takerFeeKraken,
			Exchange:               "Kraken",
			ArbitrageRecordUuid:    uuid.Nil,
			IsArbitrageOpportunity: false,
			Timestamp:              time.Now().Format(time.RFC3339),
		},
	}
	bookkeeper.RecordPriceRecord(priceRecords...)
	return priceRecords
}

// Trades returns the recent trades for given pair
func (api *KrakenAPI) Trades(pair string, since int64) (*TradesResponse, error) {
	values := url.Values{"pair": {pair}}
	if since > 0 {
		values.Set("since", strconv.FormatInt(since, 10))
	}
	resp, err := api.queryPublic("Trades", values, nil)
	if err != nil {
		return nil, err
	}

	v := resp.(map[string]interface{})

	last, err := strconv.ParseInt(v["last"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	result := &TradesResponse{
		Last:   last,
		Trades: make([]TradeInfo, 0),
	}

	trades := v[pair].([]interface{})
	for _, v := range trades {
		trade := v.([]interface{})

		priceString := trade[0].(string)
		price, _ := strconv.ParseFloat(priceString, 64)

		volumeString := trade[1].(string)
		volume, _ := strconv.ParseFloat(trade[1].(string), 64)

		tradeInfo := TradeInfo{
			Price:         priceString,
			PriceFloat:    price,
			Volume:        volumeString,
			VolumeFloat:   volume,
			Time:          int64(trade[2].(float64)),
			Buy:           trade[3].(string) == BUY,
			Sell:          trade[3].(string) == SELL,
			Market:        trade[4].(string) == MARKET,
			Limit:         trade[4].(string) == LIMIT,
			Miscellaneous: trade[5].(string),
		}

		result.Trades = append(result.Trades, tradeInfo)
	}

	return result, nil
}

// Query sends a query to Kraken api for given method and parameters
func (api *KrakenAPI) Query(method string, data map[string]string) (interface{}, error) {
	values := url.Values{}
	for key, value := range data {
		values.Set(key, value)
	}

	// Check if method is public or private
	if isStringInSlice(method, publicMethods) {
		return api.queryPublic(method, values, nil)
	} else if isStringInSlice(method, privateMethods) {
		return api.queryPrivate(method, values, nil)
	}

	return nil, fmt.Errorf("Method '%s' is not valid", method)
}

// Execute a public method query
func (api *KrakenAPI) queryPublic(method string, values url.Values, typ interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s/%s/public/%s", APIURL, APIVersion, method)
	resp, err := api.doRequest(url, values, nil, typ)

	return resp, err
}

// queryPrivate executes a private method query
func (api *KrakenAPI) queryPrivate(method string, values url.Values, typ interface{}) (interface{}, error) {
	urlPath := fmt.Sprintf("/%s/private/%s", APIVersion, method)
	reqURL := fmt.Sprintf("%s%s", APIURL, urlPath)
	secret, _ := base64.StdEncoding.DecodeString(api.secret)
	values.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	// Create signature
	signature := createSignature(urlPath, values, secret)

	// Add Key and signature to request headers
	headers := map[string]string{
		"API-Key":  api.key,
		"API-Sign": signature,
	}

	resp, err := api.doRequest(reqURL, values, headers, typ)

	return resp, err
}

// doRequest executes a HTTP Request to the Kraken API and returns the result
func (api *KrakenAPI) doRequest(reqURL string, values url.Values, headers map[string]string, typ interface{}) (interface{}, error) {

	var req *http.Request
	var err error
	if values.Get("pair") != "" {
		// Create request
		req, err = http.NewRequest("POST", reqURL, strings.NewReader(values.Encode()))
		if err != nil {
			return nil, fmt.Errorf("Could not execute request! #1 (%s)", err.Error())
		}
	} else {
		// Create request
		req, err = http.NewRequest("POST", reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("Could not execute request! #1 (%s)", err.Error())
		}
	}

	req.Header.Add("User-Agent", APIUserAgent)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Execute request
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #2 (%s)", err.Error())
	}
	defer resp.Body.Close()

	// Read request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #3 (%s)", err.Error())
	}

	// Check mime type of response
	mimeType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("Could not execute request #4! (%s)", err.Error())
	}
	if mimeType != "application/json" {
		return nil, fmt.Errorf("Could not execute request #5! (%s)", fmt.Sprintf("Response Content-Type is '%s', but should be 'application/json'.", mimeType))
	}

	// Parse request
	var jsonData KrakenResponse

	// Set the KrakenResponse.Result to typ so `json.Unmarshal` will
	// unmarshal it into given type, instead of `interface{}`.
	if typ != nil {
		jsonData.Result = typ
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #6 (%s)", err.Error())
	}

	// Check for Kraken API error
	if len(jsonData.Error) > 0 {
		return nil, fmt.Errorf("Could not execute request! #7 (%s)", jsonData.Error)
	}

	return jsonData.Result, nil
}

// isStringInSlice is a helper function to test if given term is in a list of strings
func isStringInSlice(term string, list []string) bool {
	for _, found := range list {
		if term == found {
			return true
		}
	}
	return false
}

// getSha256 creates a sha256 hash for given []byte
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func createSignature(urlPath string, values url.Values, secret []byte) string {
	// See https://www.kraken.com/help/api#general-usage for more information
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macSum := getHMacSha512(append([]byte(urlPath), shaSum...), secret)
	return base64.StdEncoding.EncodeToString(macSum)
}
