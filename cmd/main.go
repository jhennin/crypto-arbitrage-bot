package main

import (
	"cryptoArbitrageBot/arbitrageHunter"
	"cryptoArbitrageBot/internal"
	"cryptoArbitrageBot/internal/utils"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var coinbaseProURL string
var key string
var secret string
var passphrase string

/*
Initialize the environment variables.
Values are stored in a configuration file (e.g. config-DEV.yaml)
*/
func init() {
	var ok bool

	viper.SetConfigName("config-DEV")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.AutomaticEnv()

	coinbaseProURL, ok = viper.Get("COINBASE_PRO.URL").(string)
	if !ok {
		log.Fatalln("Failed to load environment variable `coinbaseProURL`. Invalid type assertion.")
	}
	key, ok = viper.Get("COINBASE_PRO.TEST.KEY").(string)
	if !ok {
		log.Fatalln("Failed to load environment variable `key`. Invalid type assertion.")
	}
	secret, ok = viper.Get("COINBASE_PRO.TEST.SECRET").(string)
	if !ok {
		log.Fatalln("Failed to load environment variable `secret`. Invalid type assertion.")
	}
	passphrase, ok = viper.Get("COINBASE_PRO.TEST.PASSPHRASE").(string)
	if !ok {
		log.Fatalln("Failed to load environment variable `passphrase`. Invalid type assertion.")
	}
	_, ok = viper.Get("KRAKEN.URL").(string)
	if !ok {
		log.Fatalf("Failed to load environment variable `KRAKEN.URL`. Invalid type assertion. %s", err.Error())
	}
	_, ok = viper.Get("DATABASE.MY_SQL_DOCKER.USERNAME").(string)
	if !ok {
		log.Fatalf("Failed to load environment variable `KRAKEN.URL`. Invalid type assertion. %s", err.Error())
	}
	utils.InitializeLogger()
	utils.SetLoggerLevel("INFO")
}

/*
Main entry point of the Crypto Arbitrage Bot.
*/
func main() {

	databaseErr := internal.ConnectToDatabase()
	if databaseErr != nil {
		utils.Logger.Error(databaseErr.Error())
	}
	arbitrageHunterError := arbitrageHunter.Start()
	if arbitrageHunterError != nil {
		utils.Logger.Error(arbitrageHunterError.Error())
	}
}
