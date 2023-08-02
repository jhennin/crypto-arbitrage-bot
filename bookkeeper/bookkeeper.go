package bookkeeper

import (
	"cryptoArbitrageBot/internal"
	"cryptoArbitrageBot/internal/utils"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PriceRecord struct {
	Uuid                   uuid.UUID `db:"uuid"`
	Timestamp              string    `db:"timestamp"`
	Currency               string    `db:"currency"`
	Price                  float64   `db:"price"`
	Fee                    float64   `db:"fee"`
	Exchange               string    `db:"exchange"`
	ArbitrageRecordUuid    uuid.UUID `db:"arbitrage_record_uuid"`
	IsArbitrageOpportunity bool      `db:"is_arbitrage_opportunity"`
}

type ArbitrageEventRecord struct {
	Uuid                   uuid.UUID `db:"uuid"`
	Timestamp              string    `db:"timestamp"`
	Currency               string    `db:"currency"`
	PriceA                 float64   `db:"price_a"`
	ExchangeA              string    `db:"exchange_a"`
	PriceB                 float64   `db:"price_b"`
	ExchangeB              string    `db:"exchange_b"`
	ProjectedProfit        float64   `db:"projected_profit"`
	IsArbitrageOpportunity bool      `db:"is_arbitrage_opportunity"`
}

type TriangularArbitrageEventRecord struct {
	Uuid                             uuid.UUID `db:"uuid"`
	Timestamp                        string    `db:"timestamp"`
	Exchange                         string    `db:"exchange"`
	TradePair1                       string    `db:"trade_pair_1"`
	TradePair1ExchangeRate           float64   `db:"trade_pair_1_exchange_rate"`
	TradePair2                       string    `db:"trade_pair_2"`
	TradePair2ExchangeRate           float64   `db:"trade_pair_2_exchange_rate"`
	TradePair3                       string    `db:"trade_pair_3"`
	TradePair3ExchangeRate           float64   `db:"trade_pair_3_exchange_rate"`
	CrossExchangeRateDiff            float64   `db:"cross_exchange_rate_difference"`
	IsTriangularArbitrageOpportunity bool      `db:"is_triangular_arbitrage_opportunity"`
}

func RecordPriceRecord(priceRecords ...PriceRecord) *internal.DatabaseError {
	database := goqu.New("mysql", internal.DbPool)
	for _, priceRecord := range priceRecords {
		aPriceRecord := goqu.Record{"uuid": priceRecord.Uuid.String(), "timestamp": priceRecord.Timestamp, "currency": priceRecord.Currency, "price": priceRecord.Price, "fee": priceRecord.Fee, "exchange": priceRecord.Exchange, "arbitrage_record_uuid": priceRecord.ArbitrageRecordUuid, "is_arbitrage_opportunity": priceRecord.IsArbitrageOpportunity}

		insertPriceRecordSQL, _, _ := database.Insert("price_records").Rows(aPriceRecord).ToSQL()

		_, err := internal.DbPool.Exec(insertPriceRecordSQL)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("Database query failed: %s", fmt.Sprintf(insertPriceRecordSQL)), zap.String("queryError,", err.Error()))
			return &internal.DatabaseError{
				Msg: err.Error(),
			}
		}
		utils.Logger.Debug("Inserted new `priceRecord` into the database.", zap.Object("priceRecord", &priceRecord))
	}

	return nil
}

func RecordArbitrageRecords(arbitrageEventRecords []ArbitrageEventRecord) error {
	database := goqu.New("mysql", internal.DbPool)

	for _, arbitrageEventRecord := range arbitrageEventRecords {
		arbitrageRecord := goqu.Record{"uuid": arbitrageEventRecord.Uuid.String(), "timestamp": arbitrageEventRecord.Timestamp, "currency": arbitrageEventRecord.Currency, "price_a": arbitrageEventRecord.PriceA, "exchange_a": arbitrageEventRecord.ExchangeA, "price_b": arbitrageEventRecord.PriceB, "exchange_b": arbitrageEventRecord.ExchangeB, "projected_profit": arbitrageEventRecord.ProjectedProfit, "is_arbitrage_opportunity": arbitrageEventRecord.IsArbitrageOpportunity}

		insertArbitrageEventSQL, _, _ := database.Insert("arbitrage_records").Rows(arbitrageRecord).ToSQL()

		_, err := internal.DbPool.Exec(insertArbitrageEventSQL)
		if err != nil {
			utils.Logger.Error(fmt.Sprintf("Database query failed: %s", fmt.Sprintf(insertArbitrageEventSQL)), zap.String("queryError,", err.Error()))
			return err
		}

		utils.Logger.Debug("Inserted new arbitrageRecord into the database.", zap.Object("arbitrageEventRecord", &arbitrageEventRecord))
	}

	return nil
}

func RecordTriangularArbitrageRecord(triangularArbitrageEventRecord TriangularArbitrageEventRecord) error {

	database := goqu.New("mysql", internal.DbPool)

	triangularArbitrageRecord := goqu.Record{"uuid": triangularArbitrageEventRecord.Uuid.String(), "timestamp": triangularArbitrageEventRecord.Timestamp, "exchange": triangularArbitrageEventRecord.Exchange, "trade_pair_1": triangularArbitrageEventRecord.TradePair1, "trade_pair_1_exchange_rate": triangularArbitrageEventRecord.TradePair1ExchangeRate, "trade_pair_2": triangularArbitrageEventRecord.TradePair2, "trade_pair_2_exchange_rate": triangularArbitrageEventRecord.TradePair2ExchangeRate, "trade_pair_3_exchange_rate": triangularArbitrageEventRecord.TradePair3ExchangeRate, "trade_pair_3": triangularArbitrageEventRecord.TradePair3, "cross_exchange_rate_difference": triangularArbitrageEventRecord.CrossExchangeRateDiff, "is_triangular_arbitrage_opportunity": triangularArbitrageEventRecord.IsTriangularArbitrageOpportunity}

	insertTriangularArbitrageEventSQL, _, _ := database.Insert("triangular_arbitrage_1_exchange").Rows(triangularArbitrageRecord).ToSQL()

	_, err := internal.DbPool.Exec(insertTriangularArbitrageEventSQL)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Database query failed: %s", fmt.Sprintf(insertTriangularArbitrageEventSQL)), zap.String("queryError,", err.Error()))
		return err
	}

	utils.Logger.Debug("Inserted new triangularArbitrageRecord into the database.", zap.Object("triangularArbitrageEventRecord", &TriangularArbitrageEventRecord{}))

	return nil
}

func (p PriceRecord) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("uuid", p.Uuid.String())
	encoder.AddString("time", p.Timestamp)
	encoder.AddString("currency", p.Currency)
	encoder.AddFloat64("price", p.Price)
	encoder.AddFloat64("fee", p.Fee)
	encoder.AddString("exchange", p.Exchange)
	encoder.AddString("arbitrage_record_uuid", p.ArbitrageRecordUuid.String())
	encoder.AddBool("is_arbitrage_opportunity", p.IsArbitrageOpportunity)
	return nil
}

func (arbitrageEventRecord ArbitrageEventRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddFloat64("priceA", arbitrageEventRecord.PriceA)
	enc.AddString("exchangeA", arbitrageEventRecord.ExchangeA)
	enc.AddFloat64("priceB", arbitrageEventRecord.PriceB)
	enc.AddString("exchangeB", arbitrageEventRecord.ExchangeB)
	enc.AddBool("isArbitrageOpportunity", arbitrageEventRecord.IsArbitrageOpportunity)
	return nil
}

func (t TriangularArbitrageEventRecord) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("uuid", t.Uuid.String())
	encoder.AddString("timestamp", t.Timestamp)
	encoder.AddString("exchange", t.Exchange)
	encoder.AddString("trade_pair_1", t.TradePair1)
	encoder.AddString("trade_pair_2", t.TradePair2)
	encoder.AddString("trade_pair_3", t.TradePair3)
	encoder.AddFloat64("cross_exchange_rate_difference", t.CrossExchangeRateDiff)
	encoder.AddBool("isTriangularArbitrageOpportunity", t.IsTriangularArbitrageOpportunity)
	return nil
}
