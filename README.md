# Crypto Arbitrage Bot

:warning: This application is not a financial advisor. I am not a financial advisor, and I do not provide any financial advice. The use of this application is at your own risk. You should do your own research before making any investment decisions. :warning:

This is a Crypto Currency Arbitrage Bot designed to detect arbitrage opportunities in the market across three different exchanges: Coinbase Pro, Gemini, and Kraken. The application deploys 2 containers via Docker: a MYSQL database, and a Grafana server. The Grafana server is used to visualize the data in the MYSQL server, as well as detecting arbitrage opportunites. It is designed to be lightweight, fast, and portable. Additionally, there is a sutie of tests that cover a number of different scenarios.

If your interested in learning more about the project, check out the article ["Run Your Own Crypto Arbitrage Bot"](medium.com).

## Pre-requisites
 * [Docker](https://docs.docker.com/get-docker/)
 * [Coinbase Pro API Key](https://help.coinbase.com/en/exchange/managing-my-account/how-to-create-an-api-key)
 * [Gemini API Key](https://support.gemini.com/hc/en-us/articles/360031080191-How-do-I-create-an-API-key-)
 * [Kraken API Key](https://support.kraken.com/hc/en-us/articles/360000919966-How-to-create-an-API-key)

1. Install Docker on your local machine
2. Create API keys for Coinbase Pro, Gemini, and Kraken
3. Update the `config-DEV.yaml` file with your API keys (see TODOs)
4. Update passwords in `config-Docker.env`, and `config-DEV.yaml` from their placeholder default values (see TODOs)

## Run the Application
1. Open Docker Desktop, or start the Docker daemon in your terminal
2. Run `docker-compose up` from the root directory of the project
3. Open a new terminal window and run `go run cmd/main.go`from the root directory of the project


## Run Tests
Run `go test` from [GoLand](https://www.jetbrains.com/go/).

Or

Run all tests in currenct directory and all sub directories:

`go test ./...`



### Resources
* [Coinbase Pro API Docs](https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getaccounts)
* [Kraken API Docs](https://docs.kraken.com/rest/)
* [Gemini API Docs](https://docs.gemini.com/rest-api/)
* [Docker network features](https://docs.docker.com/desktop/networking/)