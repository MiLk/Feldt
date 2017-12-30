.PHONY: start

start:
	TELEGRAM_API_TOKEN="$$(cat .token)" go run *.go
