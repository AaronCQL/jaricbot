package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AaronCQL/jaricbot/internal/app/bot/model"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/google/generative-ai-go/genai"
)

const (
	ActionTyping = "typing"

	ParseModeMarkdown = "Markdown"
	ParseModeHTML     = "HTML"

	ChatTypePrivate = "private"

	CommandStart = "/start"
	CommandHelp  = "/help"

	MessageStart = "Hi, I'm JaricBot! Feel free to chat about anything with me.\n\nI can continue a conversation, but only if you send your message as a reply. If you don't, I'll just assume that you're starting a new topic.\n\nIf you want me to reply you in a group, tag me at the start of your message."

	PromptJaric = "You are a Telegram bot and can only send text messages. If a user wants to continue the conversation with you, the user must reply to the message that they want to continue from. If you are prompted for your name, it's 'JaricBot', an acronym which stands for 'Just Another Rather Intelligent Chat Bot'."
)

func NewTextHandler(ctx context.Context, gen *genai.GenerativeModel, mod *model.Model) ext.Handler {
	// TODO: handle edited messages
	handler := func(bot *gotgbot.Bot, ectx *ext.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		msg := ectx.EffectiveMessage
		senderText := msg.Text
		isReply := msg.ReplyToMessage != nil
		botUsername := "@" + bot.Username

		// For messages sent in a group, only handle messages that either:
		// 1. Tagged the bot directly
		// 2. Replied to the bot's previous message
		if msg.Chat.Type != ChatTypePrivate &&
			(!strings.HasPrefix(senderText, botUsername) &&
				(!isReply || msg.ReplyToMessage.From.Id != bot.Id)) {
			return nil
		}

		// Send typing action until bot replies
		go func() {
			ectx.EffectiveChat.SendAction(bot, ActionTyping, nil)
			ticker := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-ticker.C:
					ectx.EffectiveChat.SendAction(bot, ActionTyping, nil)
				case <-ctx.Done():
					ticker.Stop()
					return
				}
			}
		}()

		// Remove the bot's username from the message
		senderText = strings.TrimSpace(strings.TrimPrefix(msg.Text, botUsername))

		// Get the full chat history
		history := model.NewHistory()
		history.Append(model.RoleUser, PromptJaric)
		if isReply {
			reply := msg.ReplyToMessage
			// Append all replies of the latest `reply`
			msgs, err := mod.GetReplies(ctx, reply.MessageId)
			if err != nil {
				return fmt.Errorf("failed to get linked messages: %v", err)
			}
			for _, m := range msgs {
				role := model.RoleUser
				if m.SenderID == bot.Id {
					role = model.RoleModel
				}
				history.Append(role, m.Content)
			}
			// Append the latest `reply`
			role := model.RoleUser
			if reply.From.Id == bot.Id {
				role = model.RoleModel
			}
			history.Append(role, reply.Text)
		}
		history.Append(model.RoleUser, senderText)

		// Call the API
		senderPart, geminiHistory := history.ToGeminiContents()
		cs := gen.StartChat()
		cs.History = geminiHistory
		res, err := cs.SendMessage(ctx, senderPart)
		if err != nil {
			return fmt.Errorf("failed to call API: %v", err)
		}

		// Reply user with chat response
		replyText, ok := (res.Candidates[0].Content.Parts[0]).(genai.Text)
		if !ok {
			return fmt.Errorf("failed to convert chat response to text")
		}
		reply, err := msg.Reply(bot,
			formatTelegramMarkdown(string(replyText)),
			&gotgbot.SendMessageOpts{
				ParseMode:             ParseModeMarkdown,
				DisableWebPagePreview: true,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		// Store messages in DB
		msgModels := []model.Message{
			*model.NewMessage(
				msg.MessageId, msg.Date, msg.Chat.Id, msg.From.Id, msg.ReplyToMessage, senderText),
			*model.NewMessage(
				reply.MessageId, reply.Date, reply.Chat.Id, reply.From.Id, reply.ReplyToMessage, reply.Text),
		}
		if err := mod.StoreMessages(ctx, msgModels); err != nil {
			return fmt.Errorf("failed to store messages into db: %v", err)
		}

		return nil
	}

	return handlers.NewMessage(message.Text, handler)
}

func NewCommandHandler(ctx context.Context) ext.Handler {
	handler := func(bot *gotgbot.Bot, ectx *ext.Context) error {
		msg := ectx.EffectiveMessage

		// Handle the `/start` and `/help` commands
		if strings.HasPrefix(msg.Text, CommandStart) ||
			strings.HasPrefix(msg.Text, CommandHelp) {
			_, err := msg.Reply(bot, MessageStart,
				&gotgbot.SendMessageOpts{ParseMode: ParseModeMarkdown},
			)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		}
		return nil
	}

	return handlers.NewMessage(message.Command, handler)
}
