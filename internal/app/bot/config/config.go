package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotKey string `json:"telegram_bot_api_key"`
	OpenAIKey      string `json:"openai_api_key"`
	PebbleDir      string `json:"pebble_dir"`
}

const (
	dotenvFile = ".env"

	KeyTelegramBotApiKey = "TELEGRAM_BOT_API_KEY"
	keyOpenaiApiKey      = "OPENAI_API_KEY"
	KeyPebbleDir         = "PEBBLE_DIR"

	DefaultPebbleDir = ".pebble"
)

func New() Config {
	log.Printf("loading %v file...\n", dotenvFile)
	godotenv.Load(dotenvFile)

	telegramBotKey, ok := os.LookupEnv(KeyTelegramBotApiKey)
	if !ok {
		panic(KeyTelegramBotApiKey + " environment variable must be set")
	}

	openaiKey, ok := os.LookupEnv(keyOpenaiApiKey)
	if !ok {
		panic(keyOpenaiApiKey + " environment variable must be set")
	}

	pebbleDir, ok := os.LookupEnv(KeyPebbleDir)
	if !ok {
		pebbleDir = DefaultPebbleDir
	}

	config := Config{
		TelegramBotKey: telegramBotKey,
		OpenAIKey:      openaiKey,
		PebbleDir:      pebbleDir,
	}

	jsonBytes, _ := json.MarshalIndent(config, "", "  ")
	log.Printf("configurations loaded:\n%v\n", string(jsonBytes))

	return config
}
