package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotKey string `json:"telegram_bot_api_key"`
	GeminiApiKey   string `json:"gemini_api_key"`
	PebbleDir      string `json:"pebble_dir"`
}

const (
	dotenvFile = ".env"

	KeyTelegramBotApiKey = "TELEGRAM_BOT_API_KEY"
	KeyGeminiApiKey      = "GEMINI_API_KEY"
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

	geminiApiKey, ok := os.LookupEnv(KeyGeminiApiKey)
	if !ok {
		panic(KeyGeminiApiKey + " environment variable must be set")
	}

	pebbleDir, ok := os.LookupEnv(KeyPebbleDir)
	if !ok {
		pebbleDir = DefaultPebbleDir
	}

	config := Config{
		TelegramBotKey: telegramBotKey,
		GeminiApiKey:   geminiApiKey,
		PebbleDir:      pebbleDir,
	}

	jsonBytes, _ := json.MarshalIndent(config, "", "  ")
	log.Printf("configurations loaded:\n%v\n", string(jsonBytes))

	return config
}
