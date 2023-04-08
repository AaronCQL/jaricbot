package model

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Message struct {
	MessageID int64  `json:"message_id"`
	CreatedAt int64  `json:"created_at"`
	ChatID    int64  `json:"chat_id"`
	SenderID  int64  `json:"sender_id"`
	ReplyID   *int64 `json:"reply_id,omitempty"`
	Content   string `json:"content"`
}

func NewMessage(
	messageID, unixTimestamp, chatID, userID int64,
	messageReplied *gotgbot.Message,
	content string,
) *Message {
	var replyID *int64
	if messageReplied != nil {
		replyID = &messageReplied.MessageId
	}
	return &Message{
		MessageID: messageID,
		CreatedAt: unixTimestamp,
		ChatID:    chatID,
		SenderID:  userID,
		ReplyID:   replyID,
		Content:   content,
	}
}
