build:
	go build -o bin/bot ./cmd/bot/main.go

test:
	go test ./... -v

dev:
	air

start:
	make build && ./bin/bot
