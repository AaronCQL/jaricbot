package main

import (
	"github.com/AaronCQL/jaricbot/internal/app/bot/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	config := config.New()

	bot, err := gotgbot.NewBot(config.TelegramBotKey, nil)
	if err != nil {
		panic(err)
	}

	_, err = bot.SetMyDescription(&gotgbot.SetMyDescriptionOpts{
		Description: "Just Another Rather Intelligent Chat Bot.",
	})
	if err != nil {
		panic(err)
	}

	_, err = bot.SetMyShortDescription(&gotgbot.SetMyShortDescriptionOpts{
		ShortDescription: "Just Another Rather Intelligent Chat Bot.",
	})
	if err != nil {
		panic(err)
	}

	_, err = bot.SetMyCommands([]gotgbot.BotCommand{
		{
			Command:     "help",
			Description: "get help",
		},
	}, nil)
	if err != nil {
		panic(err)
	}
}
