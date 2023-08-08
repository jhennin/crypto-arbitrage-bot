login-mysql-crypto-bot:
	mysql -h 127.0.0.1 -P 3307 -ucrypto_bot -p

run-bot:
	go run cmd/main.go