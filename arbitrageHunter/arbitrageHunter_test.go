package arbitrageHunter

import (
	"cryptoArbitrageBot/bookkeeper"
	"cryptoArbitrageBot/internal/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("Setup tests.")

	viper.SetConfigName("config-DEV")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	utils.InitializeLogger()
	utils.SetLoggerLevel("DEBUG")

	return func(tb testing.TB) {
		log.Println("Teardown tests.")
	}
}

func Test_isArbitrageOpportunity(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	type args struct {
		exchangePrices [][]bookkeeper.PriceRecord
	}
	tests := []struct {
		name string
		args args
		want []bookkeeper.ArbitrageEventRecord
	}{
		{
			name: "Most basic: 1 currency, 2 exchanges.",
			args: args{
				exchangePrices: [][]bookkeeper.PriceRecord{
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     25000,
							Fee:       .002,
							Exchange:  "Coinbase",
							Timestamp: "123",
						},
					},
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     26000,
							Fee:       .004,
							Exchange:  "Gemini",
							Timestamp: "123",
						},
					},
				},
			},
			want: []bookkeeper.ArbitrageEventRecord{
				{
					Currency:               "BTCUSD",
					PriceA:                 25000,
					ExchangeA:              "Coinbase",
					PriceB:                 26000,
					ExchangeB:              "Gemini",
					ProjectedProfit:        3.317647058823529,
					IsArbitrageOpportunity: true,
				},
			},
		},
		{
			name: "1 currency, 3 exchanges.",
			args: args{
				exchangePrices: [][]bookkeeper.PriceRecord{
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     23000,
							Fee:       .002,
							Exchange:  "Coinbase",
							Timestamp: "123",
						},
					},
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     24500,
							Fee:       .004,
							Exchange:  "Gemini",
							Timestamp: "123",
						},
					},
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     23050,
							Fee:       .003,
							Exchange:  "Kraken",
							Timestamp: "123",
						},
					},
				},
			},
			want: []bookkeeper.ArbitrageEventRecord{
				{
					Currency:               "BTCUSD",
					PriceA:                 23000,
					ExchangeA:              "Coinbase",
					PriceB:                 24500,
					ExchangeB:              "Gemini",
					ProjectedProfit:        5.7094736842105265,
					IsArbitrageOpportunity: true,
				},
				{
					Currency:               "BTCUSD",
					PriceA:                 23000,
					ExchangeA:              "Coinbase",
					PriceB:                 23050,
					ExchangeB:              "Kraken",
					ProjectedProfit:        -0.28295331161780674,
					IsArbitrageOpportunity: false,
				},
				{
					Currency:               "BTCUSD",
					PriceA:                 24500,
					ExchangeA:              "Gemini",
					PriceB:                 23050,
					ExchangeB:              "Kraken",
					ProjectedProfit:        5.395793901156677,
					IsArbitrageOpportunity: true,
				},
			},
		},
		{
			name: "2 currency, 3 exchanges.",
			args: args{
				exchangePrices: [][]bookkeeper.PriceRecord{
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     25000,
							Fee:       .002,
							Exchange:  "Coinbase",
							Timestamp: "123",
						},
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "ETHUSD",
							Price:     1500,
							Fee:       .002,
							Exchange:  "Coinbase",
							Timestamp: "123",
						},
					},
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     23900,
							Fee:       .004,
							Exchange:  "Gemini",
							Timestamp: "123",
						},
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "ETHUSD",
							Price:     1650,
							Fee:       .004,
							Exchange:  "Gemini",
							Timestamp: "123",
						},
					},
					{
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "BTCUSD",
							Price:     25980,
							Fee:       .003,
							Exchange:  "Kraken",
							Timestamp: "123",
						},
						bookkeeper.PriceRecord{
							Uuid:      uuid.New(),
							Currency:  "ETHUSD",
							Price:     1645,
							Fee:       .003,
							Exchange:  "Kraken",
							Timestamp: "123",
						},
					},
				},
			},
			want: []bookkeeper.ArbitrageEventRecord{
				{
					Currency:               "BTCUSD",
					PriceA:                 25000,
					ExchangeA:              "Coinbase",
					PriceB:                 23900,
					ExchangeB:              "Gemini",
					ProjectedProfit:        3.903476482617587,
					IsArbitrageOpportunity: true,
				},
				{
					Currency:               "BTCUSD",
					PriceA:                 25000,
					ExchangeA:              "Coinbase",
					PriceB:                 25980,
					ExchangeB:              "Kraken",
					ProjectedProfit:        3.342722636327972,
					IsArbitrageOpportunity: true,
				},
				{
					Currency:               "BTCUSD",
					PriceA:                 23900,
					ExchangeA:              "Gemini",
					PriceB:                 25980,
					ExchangeB:              "Kraken",
					ProjectedProfit:        7.644186046511628,
					IsArbitrageOpportunity: true,
				},
				{
					Currency:               "ETHUSD",
					PriceA:                 1500,
					ExchangeA:              "Coinbase",
					PriceB:                 1650,
					ExchangeB:              "Gemini",
					ProjectedProfit:        8.914285714285715,
					IsArbitrageOpportunity: true,
				},
				{
					Currency:               "ETHUSD",
					PriceA:                 1500,
					ExchangeA:              "Coinbase",
					PriceB:                 1645,
					ExchangeB:              "Kraken",
					ProjectedProfit:        8.716375198728139,
					IsArbitrageOpportunity: true,
				},
				{
					Currency:               "ETHUSD",
					PriceA:                 1650,
					ExchangeA:              "Gemini",
					PriceB:                 1645,
					ExchangeB:              "Kraken",
					ProjectedProfit:        -0.3966616084977238,
					IsArbitrageOpportunity: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isArbitrageOpportunity(tt.args.exchangePrices...)
			for i := range tt.want {
				assert.Equal(t, got[i].PriceA, tt.want[i].PriceA)
				assert.Equal(t, got[i].PriceB, tt.want[i].PriceB)
				assert.Equal(t, got[i].ExchangeA, tt.want[i].ExchangeA)
				assert.Equal(t, got[i].ExchangeB, tt.want[i].ExchangeB)
				assert.Equal(t, got[i].ProjectedProfit, tt.want[i].ProjectedProfit)
				assert.Equal(t, got[i].IsArbitrageOpportunity, tt.want[i].IsArbitrageOpportunity)
			}
		})
	}
}

func Test_isTrangularArbitrage1Exchange(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	type args struct {
		exchangePrices []bookkeeper.PriceRecord
	}
	tests := []struct {
		name string
		args args
		want bookkeeper.TriangularArbitrageEventRecord
	}{
		{
			name: "Is triangular arbitrage opportunity (ETHBTC, LTCETH, LTCBTC); NO FEE",
			args: args{
				exchangePrices: []bookkeeper.PriceRecord{
					bookkeeper.PriceRecord{
						Uuid:      uuid.New(),
						Currency:  "ETHBTC",
						Price:     .06,
						Fee:       0,
						Exchange:  "Gemini",
						Timestamp: "123",
					},
					bookkeeper.PriceRecord{
						Uuid:      uuid.New(),
						Currency:  "LTCETH",
						Price:     .04776,
						Fee:       0,
						Exchange:  "Gemini",
						Timestamp: "123",
					},
					bookkeeper.PriceRecord{
						Uuid:      uuid.New(),
						Currency:  "LTCBTC",
						Price:     0.003144,
						Fee:       0,
						Exchange:  "Gemini",
						Timestamp: "123",
					},
				},
			},
			want: bookkeeper.TriangularArbitrageEventRecord{
				Uuid:                             uuid.UUID{},
				Timestamp:                        "2023-07-26T18:00:00-00:00",
				Exchange:                         "Gemini",
				TradePair1:                       "ETHBTC",
				TradePair1ExchangeRate:           .06,
				TradePair2:                       "LTCETH",
				TradePair2ExchangeRate:           20.938023450586265,
				TradePair3:                       "LTCBTC",
				TradePair3ExchangeRate:           0.003144,
				CrossExchangeRateDiff:            0.09715242881072039,
				IsTriangularArbitrageOpportunity: true,
			},
		},
		{
			name: "NO triangular arbitrage opportunity (ETHBTC, LTCETH, LTCBTC)",
			args: args{
				exchangePrices: []bookkeeper.PriceRecord{
					bookkeeper.PriceRecord{
						Uuid:      uuid.New(),
						Currency:  "ETHBTC",
						Price:     .06,
						Fee:       .002,
						Exchange:  "Gemini",
						Timestamp: "123",
					},
					bookkeeper.PriceRecord{
						Uuid:      uuid.New(),
						Currency:  "LTCETH",
						Price:     .07,
						Fee:       .002,
						Exchange:  "Gemini",
						Timestamp: "123",
					},
					bookkeeper.PriceRecord{
						Uuid:      uuid.New(),
						Currency:  "LTCBTC",
						Price:     0.003,
						Fee:       .002,
						Exchange:  "Gemini",
						Timestamp: "123",
					},
				},
			},
			want: bookkeeper.TriangularArbitrageEventRecord{
				Uuid:                             uuid.UUID{},
				Timestamp:                        "2023-07-26T18:00:00-00:00",
				Exchange:                         "Gemini",
				TradePair1:                       "ETHBTC",
				TradePair1ExchangeRate:           .06,
				TradePair2:                       "LTCETH",
				TradePair2ExchangeRate:           14.285714285714285,
				TradePair3:                       "LTCBTC",
				TradePair3ExchangeRate:           0.003,
				CrossExchangeRateDiff:            -0.2917142857142857,
				IsTriangularArbitrageOpportunity: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTrangularArbitrage1Exchange(tt.args.exchangePrices)
			assert.Equal(t, got.Exchange, tt.want.Exchange)
			assert.Equal(t, got.TradePair1, tt.want.TradePair1)
			assert.Equal(t, got.TradePair1ExchangeRate, tt.want.TradePair1ExchangeRate)
			assert.Equal(t, got.TradePair2, tt.want.TradePair2)
			assert.Equal(t, got.TradePair2ExchangeRate, tt.want.TradePair2ExchangeRate)
			assert.Equal(t, got.TradePair3, tt.want.TradePair3)
			assert.Equal(t, got.TradePair3ExchangeRate, tt.want.TradePair3ExchangeRate)
			assert.Equal(t, got.CrossExchangeRateDiff, tt.want.CrossExchangeRateDiff)
			assert.Equal(t, got.IsTriangularArbitrageOpportunity, tt.want.IsTriangularArbitrageOpportunity)
		})
	}
}
