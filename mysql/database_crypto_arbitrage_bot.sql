USE crypto_arbitrage_bot;

CREATE TABLE `arbitrage_records` (
  `uuid` varchar(255) NOT NULL,
  `timestamp` timestamp NOT NULL,
  `currency` varchar(255) NOT NULL,
  `exchange_a` varchar(255) NOT NULL,
  `price_a` varchar(255) NOT NULL,
  `exchange_b` varchar(255) NOT NULL,
  `price_b` varchar(255) NOT NULL,
  `projected_profit` double NOT NULL,
  `is_arbitrage_opportunity` varchar(255) NOT NULL,
  PRIMARY KEY (`uuid`)
);

CREATE TABLE `triangular_arbitrage_1_exchange` (
    `uuid` varchar(255) NOT NULL,
    `timestamp` timestamp NOT NULL,
    `exchange` varchar(255) NOT NULL,
    `trade_pair_1` varchar(255) NOT NULL,
    `trade_pair_1_exchange_rate` double NOT NULL,
    `trade_pair_2` varchar(255) NOT NULL,
    `trade_pair_2_exchange_rate` double NOT NULL,
    `trade_pair_3` varchar(255) NOT NULL,
    `trade_pair_3_exchange_rate` double NOT NULL,
    `cross_exchange_rate_difference` double NOT NULL,
    `is_triangular_arbitrage_opportunity` varchar(255) NOT NULL,
    PRIMARY KEY (`uuid`)
);

CREATE TABLE `price_records` (
    `uuid` varchar(255) NOT NULL,
    `timestamp` timestamp NOT NULL,
    `currency` varchar(255) NOT NULL,
    `price` double NOT NULL,
    `fee` varchar(255) NOT NULL,
    `exchange` varchar(255) NOT NULL,
    `arbitrage_record_uuid` varchar(255) NOT NULL,
    `is_arbitrage_opportunity` varchar(255) NOT NULL,
    PRIMARY KEY (`uuid`)
);


-- crypto_bot MySQL user
CREATE USER 'crypto_bot'@'%' IDENTIFIED BY 'change_crypto_bot_password';
REVOKE USAGE ON *.* FROM 'crypto_bot'@'%';
GRANT SELECT, INSERT, DELETE ON crypto_arbitrage_bot.* TO 'crypto_bot'@'%';

-- Grafana MySQL user
REVOKE ALL PRIVILEGES ON *.* FROM 'grafana'@'%';
GRANT SELECT ON arbitrage_records TO 'grafana'@'%';
GRANT SELECT ON price_records TO 'grafana'@'%';
GRANT SELECT ON triangular_arbitrage_1_exchange TO 'grafana'@'%';

FLUSH PRIVILEGES;
