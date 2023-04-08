build:
	go build -o bin/bot ./cmd/bot/main.go

test:
	go test ./... -v

setup:
	go run ./cmd/setup/main.go

dev:
	air

start:
	make build && ./bin/bot
