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
	"github.com/sashabaranov/go-openai"
)

const (
	ActionTyping = "typing"

	ParseModeMarkdown = "Markdown"
	ParseModeHTML     = "HTML"

	ChatTypePrivate = "private"

	CommandStart = "/start"
	CommandHelp  = "/help"

	MessageStart = "Hi, I'm JaricBot! Feel free to chat about anything with me.\n\nI can continue a conversation, but only if you send your message as a reply. If you don't, I'll just assume that you're starting a new topic.\n\nI can also reply you in a group, but you'll need to tag me at the start of your message."

	PromptJaric = "Your name is 'JaricBot', an acronym which stands for 'Just Another Rather Intelligent Chat Bot'."
)

func NewTextHandler(ctx context.Context, client *openai.Client, mod *model.Model) ext.Handler {
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

		// Handle prompts and replies
		ccMsgs := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: PromptJaric,
			},
		}
		if isReply {
			reply := msg.ReplyToMessage
			// Append all replies of the latest `reply`
			msgs, err := mod.GetReplies(ctx, reply.MessageId)
			if err != nil {
				return fmt.Errorf("failed to get linked messages: %v", err)
			}
			for _, m := range msgs {
				role := openai.ChatMessageRoleUser
				if m.SenderID == bot.Id {
					role = openai.ChatMessageRoleAssistant
				}
				ccMsgs = append(ccMsgs, openai.ChatCompletionMessage{
					Role:    role,
					Content: m.Content,
				})
			}
			// Append the latest `reply`
			role := openai.ChatMessageRoleUser
			if reply.From.Id == bot.Id {
				role = openai.ChatMessageRoleAssistant
			}
			ccMsgs = append(ccMsgs, openai.ChatCompletionMessage{
				Role:    role,
				Content: reply.Text,
			})
		}
		ccMsgs = append(ccMsgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: senderText,
		})

		// Query OpenAI to get chat completion
		res, err := client.CreateChatCompletion(ctx,
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: ccMsgs,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create chat completion: %v", err)
		}

		// Reply user with chat response
		reply, err := msg.Reply(bot,
			formatTelegramMarkdown(res.Choices[0].Message.Content),
			&gotgbot.SendMessageOpts{
				ParseMode: ParseModeMarkdown,
				ReplyMarkup: &gotgbot.ForceReply{
					ForceReply: true,
				},
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
