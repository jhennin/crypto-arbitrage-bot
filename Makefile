create-arbitrage-bot-databases:
	docker exec -it mysql_docker mysql -uroot -p -e "CREATE DATABASE crypto_arbitrage_bot;"

create-grafana-db-user:
	docker exec -it mysql_docker mysql -uroot -p -e "CREATE USER 'grafana'@'localhost' IDENTIFIED BY 'test';"

give-grafana-db-user-permissions:
	docker exec -it mysql_docker mysql -uroot -p -e "GRANT ALL PRIVILEGES ON crypto_arbitrage_bot.arbitrage_records TO 'grafana'@'localhost';"

login-mysql-docker:
	docker exec -it mysql_docker mysql -uroot -p;

create-arbitrage-events-table:
	docker exec -it mysql_docker mysql -uroot -ptest -e "CREATE TABLE crypto_arbitrage_bot.arbitrage_records (uuid VARCHAR(255) NOT NULL, PRIMARY KEY (uuid), timestamp VARCHAR(255), price_a VARCHAR(255), exchange_a VARCHAR(255), price_b VARCHAR(255), exchange_b VARCHAR(255), is_arbitrage_opportunity VARCHAR(255));"

create-price-records-table:
	docker exec -it mysql_docker mysql -uroot -ptest -e "CREATE TABLE crypto_arbitrage_bot.price_records (uuid VARCHAR(255) NOT NULL, PRIMARY KEY (uuid), timestamp VARCHAR(255), currency VARCHAR(255), price VARCHAR(255), fee VARCHAR(255), exchange VARCHAR(255), arbitrage_record_id VARCHAR(255), is_arbitrage_opportunity VARCHAR(255));"

# --> Log into mysql from local host computer terminal
	# mysql -h 127.0.0.1 -P 3307 -uroot -p


######################################################
######			USEFUL MYSQL COMMANDS			######
######################################################

# --> Change the database password
# ALTER USER 'root'@'localhost' IDENTIFIED BY 'test';

# --> Create a new test table
# CREATE TABLE test.test (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id), name VARCHAR(255), age INT);


# --> Find what IP address docker container is at
# docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql_docker

# --> Create a record in the database
# INSERT INTO crypto_arbitrage_bot.arbitrage_bot (exchange, pair, buy_price, sell_price, profit, profit_percent, timestamp) VALUES ('Binance', 'BTC/USDT', '10000', '10001', '1', '0.01', '2021-01-01 00:00:00');

# --> Delete table in database
# DROP TABLE crypto_arbitrage_bot.price_record;

# --> Show tables in the database
# SHOW TABLES FROM crypto_arbitrage_bot;

# --> Delete all records in the table
# DELETE FROM crypto_arbitrage_bot.price_records;

# --> Show all records in the table
# SELECT * FROM crypto_arbitrage_bot.arbitrage_bot;

# --> Show columns in the table
# SHOW COLUMNS FROM crypto_arbitrage_bot.arbitrage_records;

# Docker network features
# https://docs.docker.com/desktop/networking/
