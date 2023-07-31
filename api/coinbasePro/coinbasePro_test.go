package coinbasePro

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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

	viper.SetConfigName("config-TEST")
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

func Test_RequestBuilder(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	now := strconv.FormatInt(time.Now().Unix(), 10)

	key, ok := viper.Get("COINBASE_PRO.TEST.KEY").(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	secret, ok := viper.Get("COINBASE_PRO.TEST.SECRET").(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	passphrase, ok := viper.Get("COINBASE_PRO.TEST.PASSPHRASE").(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	prehashString := now + "GET" + "/test" + ""
	hmacKey, _ := base64.StdEncoding.DecodeString(secret)
	signature := hmac.New(sha256.New, hmacKey)
	signature.Write([]byte(prehashString))
	cbAccessSignature := base64.StdEncoding.EncodeToString(signature.Sum(nil))

	actualRequest := requestBuilder(now, "GET", "/test", "")

	assert.Equal(t, "application/json", actualRequest.Header.Get("Content-Type"), "FAILED: Content_Type")
	assert.Equal(t, key, actualRequest.Header.Get("CB-ACCESS-KEY"), "FAILED: CB-ACCESS-KEY")
	assert.Equal(t, cbAccessSignature, actualRequest.Header.Get("CB-ACCESS-SIGN"), "FAILED: CB-ACCESS-SIGN")
	assert.Equal(t, now, actualRequest.Header.Get("CB-ACCESS-TIMESTAMP"), "FAILED: CB-ACCESS-TIMESTAMP")
	assert.Equal(t, passphrase, actualRequest.Header.Get("CB-ACCESS-PASSPHRASE"), "FAILED: CB-ACCESS-PASSPHRASE")
}

func Test_getSignedPrices(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	client := &http.Client{Timeout: time.Second * 10}

	actualSigendPrices, err := getSignedPrices(client)

	assert2.Emptyf(t, err, "Failed to retrieve signed prices.")
	assert2.NotEmpty(t, actualSigendPrices, "Failed to retrieve signed prices.")
}
