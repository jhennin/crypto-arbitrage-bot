package gemini

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	assert2 "github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("Setup tests.")

	viper.SetConfigName("config-DEV")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return func(tb testing.TB) {
		log.Println("Teardown tests.")
	}
}

func TestRequestBuilder(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	path := "/v1/pricefeed"
	method := "GET"

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

	actualRequest := requestBuilder(now, path, method)

	assert.Equal(t, actualRequest.Header.Get("Content-Length"), "0", "FAILED: Content-Length")
	assert.Equal(t, actualRequest.Header.Get("Content-Type"), "text/plain", "FAILED: Content-Type")
	assert.Equal(t, actualRequest.Header.Get("X-GEMINI-APIKEY"), key, "FAILED: X-GEMINI-APIKEY")
	assert.Equal(t, actualRequest.Header.Get("X-GEMINI-PAYLOAD"), base64.StdEncoding.EncodeToString(encodedPayload), "FAILED: X-GEMINI-PAYLOAD")
	assert.Equal(t, actualRequest.Header.Get("X-GEMINI-SIGNATURE"), xGeminiSignature, "FAILED: X-GEMINI-SIGNATURE")
	assert.Equal(t, actualRequest.Header.Get("Cache-Control"), "no-cache", "FAILED: Cache-Control")
}

func TestGetPriceFeed(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	client := &http.Client{Timeout: time.Second * 10}

	actualRes, err := getPriceFeed(client)

	assert2.Emptyf(t, err, "Failed to retrieve price feed.")
	assert2.NotEmpty(t, actualRes, "Failed to retrieve price feed.")
}
