package arbitrageHunter

import (
	"cryptoArbitrageBot/api/coinbasePro"
	"cryptoArbitrageBot/api/gemini"
	"cryptoArbitrageBot/api/kraken"
	"cryptoArbitrageBot/bookkeeper"
	"cryptoArbitrageBot/internal"
	"cryptoArbitrageBot/internal/utils"
	"fmt"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"math"
	"sort"
	"strconv"
	"time"
)

var arbirtageHunterOutputArt = "                                             |\n                                                        \\.\n                                                        /|.\n                                                      /  `|.\n                                                    /     |.\n                                                  /       |.\n                                                /         `|.\n                                              /            |.\n                                            /              |.\n                                          /                |.\n     __                                 /                  `|.\n      -\\                              /                     |.\n        \\\\                          /                       |.\n          \\\\                      /                         |.\n           \\|                   /                           |\\\n             \\#####\\          /                             ||\n         ==###########>     /                               ||\n          \\##==      \\    /                                 ||\n     ______ =       =|__/___                                ||\n ,--' ,----`-,__ ___/'  --,-`-==============================##==========>\n\\               '        ##_______ ______   ______,--,____,=##,__\n `,    __==    ___,-,__,--'#'  ==='      `-'              | ##,-/\n   `-,____,---'       \\####\\              |        ____,--\\_##,/\n       #_              |##   \\  _____,---==,__,---'         ##\n        #              ]===--==\\                            ||\n        #,             ]         \\                          ||\n         #_            |           \\                        ||\n          ##_       __/'             \\                      ||\n           ####='     |                \\                    |/\n            ###       |                  \\                  |.\n            ##       _'                    \\                |.\n           ###=======]                       \\              |.\n          ///        |                         \\           ,|.\n          //         |                           \\         |.\n                                                   \\      ,|.\n                                                     \\    |.\n                                                       \\  |.\n                                                         \\|.\n                                                         /.\n                                                        |\n\n\n"

type ArbitrageHunterError struct {
	msg string
}

func (e *ArbitrageHunterError) Error() string {
	return fmt.Sprint("An error has occured with the Arbitrage Hunter.")
}

func Start() *ArbitrageHunterError {
	if arbitrageHunterPrompt() == false {
		return &ArbitrageHunterError{
			msg: "Unable to start the Arbitrage Hunter.",
		}
	}

	coinbaseProClient := coinbasePro.NewClient()
	geminiClient := gemini.NewClient()
	krakenCLient := kraken.NewWithClient(viper.Get("KRAKEN.API_KEY_1675894563504.KEY").(string), viper.Get("KRAKEN.API_KEY_1675894563504.KEY").(string), internal.GetClient())

	scheduler := gocron.NewScheduler(time.UTC)
	job, err := scheduler.Every(5).Seconds().Do(
		func() {
			coinbasePriceRecords, coinbaseProErr := coinbaseProClient.GetPrices()
			if coinbaseProErr != nil {
				utils.Logger.Error(fmt.Sprintf("Error fetching coinbase pro prices: %v", coinbaseProErr))
			}
			geminiPriceRecords := geminiClient.GetPrices()
			krakenPriceRecords := krakenCLient.GetPrices()
			bookkeeper.RecordArbitrageRecords(isArbitrageOpportunity(coinbasePriceRecords, geminiPriceRecords, krakenPriceRecords)) // TODO need to handle what happens when one of these records is nil. I believe this is the source of the main bug causing the program to crash.
			bookkeeper.RecordTriangularArbitrageRecord(isTrangularArbitrage1Exchange(geminiPriceRecords))
			utils.Logger.Info("Ran arbitrage hunter job.")
		})
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error running job: %v, err: %v", job, err))
	}
	scheduler.StartBlocking()

	utils.Logger.Info(fmt.Sprintf("Job: %v, err: %v", job, err))

	return nil
}

func arbitrageHunterPrompt() bool {
	fmt.Printf("Grafana dashboard is availabe at --> http://localhost:3000\n")
	fmt.Printf("Would you like to release the Arbitrage Hunter? (y/n): ")
	var startArbitrageHunter string

	fmt.Scanln(&startArbitrageHunter)
	if startArbitrageHunter != "y" {
		return false
	}

	utils.Logger.Info("===> Releasing the Arbitrage Hunter in 10 seconds...")
	fmt.Printf(arbirtageHunterOutputArt)

	for i := 10; i > 0; i-- {
		fmt.Printf("%v...", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n")

	return true
}

func isArbitrageOpportunity(exchangePrices ...[]bookkeeper.PriceRecord) []bookkeeper.ArbitrageEventRecord {
	priceRecords := flatten(exchangePrices)
	utils.Logger.Debug("------------> exchangePrices length = " + strconv.Itoa(len(priceRecords)))
	sort.Slice(priceRecords, func(i, j int) bool {
		return priceRecords[i].Currency < priceRecords[j].Currency
	})
	utils.Logger.Debug(fmt.Sprintf("Done sorting the exchange price records. Total count = %v", strconv.Itoa(len(priceRecords))))

	var arbitrageRecords []bookkeeper.ArbitrageEventRecord
	for i := 0; i < len(priceRecords)-1; i++ {
		for j := i + 1; j < len(priceRecords); j++ {
			utils.Logger.Debug(fmt.Sprintf("Current counters for the arbitrageRecords slice --> i = %v; j = %v", i, j))
			if priceRecords[i].Currency != priceRecords[j].Currency {
				continue
			}
			utils.Logger.Debug(fmt.Sprintf("Found a match for the currencies, for the following exchanges; %v, %v. Currency match is --> %v & %v", priceRecords[i].Exchange, priceRecords[j].Exchange, priceRecords[i].Currency, priceRecords[j].Currency))

			projectedProfit := ((math.Abs(priceRecords[i].Price-priceRecords[j].Price) - ((priceRecords[i].Price * priceRecords[i].Fee) + (priceRecords[j].Price * priceRecords[j].Fee))) / ((priceRecords[i].Price + priceRecords[j].Price) / 2)) * 100
			record := bookkeeper.ArbitrageEventRecord{
				Uuid:                   uuid.New(),
				Timestamp:              time.Now().Format(time.RFC3339),
				Currency:               priceRecords[i].Currency,
				PriceA:                 priceRecords[i].Price,
				ExchangeA:              priceRecords[i].Exchange,
				PriceB:                 priceRecords[j].Price,
				ExchangeB:              priceRecords[j].Exchange,
				ProjectedProfit:        projectedProfit,
				IsArbitrageOpportunity: false,
			}
			if projectedProfit > 0 {
				record.ProjectedProfit = projectedProfit
				record.IsArbitrageOpportunity = true
				utils.Logger.Info(fmt.Sprintf("Found an arbitrage opportunity! Projected profit = %v", projectedProfit))
			}
			arbitrageRecords = append(arbitrageRecords, record)
		}
	}
	return arbitrageRecords
}

func isTrangularArbitrage1Exchange(exchangePrices []bookkeeper.PriceRecord) bookkeeper.TriangularArbitrageEventRecord {
	var ethBtcExchangeRate float64
	var ethLtcExchangeRate float64 //NOTE this is inverted (i.e. ETHLTC = 1 / LTCETH)
	var ltcbtcExchangeRate float64

	for _, priceRecord := range exchangePrices {
		if priceRecord.Currency == "ETHBTC" {
			ethBtcExchangeRate = priceRecord.Price
		} else if priceRecord.Currency == "LTCETH" {
			ethLtcExchangeRate = 1 / priceRecord.Price //NOTE this is inverted (i.e. ETHLTC = 1 / LTCETH)
		} else if priceRecord.Currency == "LTCBTC" {
			ltcbtcExchangeRate = priceRecord.Price
		}
	}
	crossExchangeRate := ethLtcExchangeRate * ltcbtcExchangeRate

	crossExchangeRateDiff := ((crossExchangeRate - ethBtcExchangeRate) / ethBtcExchangeRate) - (exchangePrices[0].Fee * 3)
	isTriangularArbitrageOpportunity := func() bool {
		if crossExchangeRateDiff > 0 { // TODO There are some edge cases where crossExchangeRateDiff > 0, but isTriangularArbitrageOpportunity is false. Need to figure out why.
			utils.Logger.Info(fmt.Sprintf("Found a triangular arbitrage opportunity! CrossExchangeRateDiff = %v", crossExchangeRateDiff))
			return true
		}
		return false
	}

	record := bookkeeper.TriangularArbitrageEventRecord{
		Uuid:                             uuid.New(),
		Timestamp:                        time.Now().Format(time.RFC3339),
		Exchange:                         exchangePrices[0].Exchange,
		TradePair1:                       "ETHBTC",
		TradePair1ExchangeRate:           ethBtcExchangeRate,
		TradePair2:                       "LTCETH",
		TradePair2ExchangeRate:           ethLtcExchangeRate,
		TradePair3:                       "LTCBTC",
		TradePair3ExchangeRate:           ltcbtcExchangeRate,
		CrossExchangeRateDiff:            crossExchangeRateDiff,
		IsTriangularArbitrageOpportunity: isTriangularArbitrageOpportunity(),
	}

	return record
}

func flatten(m [][]bookkeeper.PriceRecord) []bookkeeper.PriceRecord {

	var priceRecords []bookkeeper.PriceRecord
	for _, recordSet := range m {
		priceRecords = append(priceRecords, recordSet...)
	}
	utils.Logger.Debug(fmt.Sprintf("totalRecords = %v", strconv.Itoa(len(priceRecords))))
	return priceRecords
}
