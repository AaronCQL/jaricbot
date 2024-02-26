package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/AaronCQL/jaricbot/internal/app/bot/config"
	"github.com/AaronCQL/jaricbot/internal/app/bot/database"
	"github.com/AaronCQL/jaricbot/internal/app/bot/handler"
	"github.com/AaronCQL/jaricbot/internal/app/bot/model"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	storeMessages string = "messages"
	storeUsers    string = "users"

	maxUpdates      int64 = 8  // Limit number of concurrent updates
	pollingInterval int64 = 16 // Max that telegram allows is 50s
)

func main() {
	config := config.New()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	db := database.New(config.PebbleDir, storeMessages)
	defer db.Close()
	model := model.New(db)
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GeminiApiKey))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	gen := client.GenerativeModel("gemini-pro")
	gen.SetCandidateCount(1)
	gen.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// Log errors returned by handlers
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println(err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})
	dispatcher := updater.Dispatcher

	// Add handlers here:
	dispatcher.AddHandler(handler.NewCommandHandler(ctx))
	dispatcher.AddHandler(handler.NewTextHandler(ctx, gen, model))

	bot, err := gotgbot.NewBot(config.TelegramBotKey, nil)
	if err != nil {
		log.Printf("failed to initialise telegram bot\n")
		panic(err)
	}
	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: false,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Limit:   maxUpdates,
			Timeout: pollingInterval,
			RequestOpts: &gotgbot.RequestOpts{
				// This timeout must be longer than the polling interval
				Timeout: time.Second * (time.Duration(pollingInterval) + 1),
			},
			AllowedUpdates: []string{
				gotgbot.UpdateTypeMessage,
				gotgbot.UpdateTypeEditedMessage,
			},
		},
	})
	if err != nil {
		log.Printf("failed to start telegram bot\n")
		panic(err)
	}

	log.Printf("telegram bot [%v] started!\n", bot.Username)

	<-ctx.Done()

	log.Printf("exiting gracefully...\n")
}
