package handler

import (
	"context"
	"fmt"
	"strings"

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

	CommandPrefix = "/"
	CommandStart  = "/start"
	CommandHelp   = "/help"
	CommandChat   = "/chat"

	MessageStart = "Hi, I'm JaricBot! Feel free to chat about anything with me.\n\nJust take note: if you want to continue the conversation, please *send your message as a reply*."

	PromptJaric = "Your name is 'JaricBot', an acronym which stands for 'Just Another Rather Intelligent Chat Bot'."
)

func NewTextHandler(ctx context.Context, client *openai.Client, mod *model.Model) ext.Handler {
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

		// Handle the `/chat` command
		senderText := strings.TrimSpace(strings.TrimPrefix(msg.Text, CommandChat))
		// Ignore all other commands
		if strings.HasPrefix(senderText, CommandPrefix) {
			return nil
		}

		// TODO: handle edited messages

		// Send typing action
		go ectx.EffectiveChat.SendAction(bot, ActionTyping, nil)

		// Query past messages if user is replying to a message
		ccMsgs := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: PromptJaric,
			},
		}
		if msg.ReplyToMessage != nil {
			msgs, err := mod.GetMessageAndReplies(ctx, msg.ReplyToMessage.MessageId)
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
			&gotgbot.SendMessageOpts{ParseMode: ParseModeMarkdown},
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

	return handlers.NewMessage(message.All, handler)
}
