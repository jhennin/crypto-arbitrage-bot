package internal

import (
	"context"
	"cryptoArbitrageBot/internal/utils"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
	"time"
)

var (
	DbPool *sql.DB
	err    error
)

type DatabaseError struct {
	Msg string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprint("Error establishing connection to database: ") + e.Msg
}

/*
Connect to database and test connection. Returns error if connection fails.
*/
func ConnectToDatabase() *DatabaseError {

	DbPool, err = connectTCPSocket()
	if err != nil {
		return &DatabaseError{
			Msg: err.Error(),
		}
	}
	DbPool.SetMaxIdleConns(0)
	DbPool.SetConnMaxLifetime(time.Second)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if dbErr := DbPool.PingContext(ctx); dbErr != nil {
		return &DatabaseError{
			Msg: dbErr.Error(),
		}
	}
	err := testDatabase()
	if err != nil {
		return err
	}
	utils.Logger.Info("Database is online and fully operational. All database tests passed.")
	return nil
}

/*
Initializes a TCP connection pool for a Cloud SQL instance of MySQL.
*/
func connectTCPSocket() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := viper.Get(k).(string)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}

	var (
		dbUser    = mustGetenv("DATABASE.MY_SQL_DOCKER.USERNAME")
		dbPwd     = mustGetenv("DATABASE.MY_SQL_DOCKER.PASSWORD")
		dbName    = mustGetenv("DATABASE.MY_SQL_DOCKER.NAME")
		dbTCPHost = mustGetenv("DATABASE.MY_SQL_DOCKER.HOST")
		dbPort    = mustGetenv("DATABASE.MY_SQL_DOCKER.PORT")
	)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	utils.Logger.Debug(dbURI)
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
}

/*
Tests the database by inserting a test record and then deleting it.
*/
func testDatabase() *DatabaseError {
	database := goqu.New("mysql", DbPool)

	arbitrageTestRecord := goqu.Record{"uuid": "1234", "timestamp": "2023-02-17 14:41:00", "currency": "BTCUSD", "exchange_a": "Binance", "price_a": "1000.00", "exchange_b": "Binance", "price_b": "1050.00", "projected_profit": "3.33", "is_arbitrage_opportunity": "0"}

	insertArbitrageEventSQL, _, _ := database.Insert("arbitrage_records").Rows(arbitrageTestRecord).ToSQL()
	removeArbitrageEventSQL, _, _ := database.Delete("arbitrage_records").Where(goqu.Ex{"uuid": "1234"}).ToSQL()

	_, err := DbPool.Exec(insertArbitrageEventSQL)
	if err != nil {
		return &DatabaseError{
			Msg: fmt.Sprintf("Database test query failed: %s. %v", fmt.Sprintf(insertArbitrageEventSQL), err),
		}
	}
	utils.Logger.Info("(TEST) inserted test arbitrage event into the database.")

	_, err = DbPool.Exec(removeArbitrageEventSQL)
	if err != nil {
		return &DatabaseError{
			Msg: fmt.Sprintf("Database test query failed: %s. %v", fmt.Sprintf(removeArbitrageEventSQL), err),
		}
	}
	utils.Logger.Info("(TEST) deleted test arbitrage event into the database.")
	return nil
}
