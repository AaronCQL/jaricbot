package model

import (
	"context"

	"github.com/AaronCQL/jaricbot/internal/app/bot/database"
	"github.com/cockroachdb/pebble"
)

type Model struct {
	db *database.Database
}

const (
	maxLinkedMessages int = 10
)

func New(db *database.Database) *Model {
	return &Model{
		db: db,
	}
}

// Stores messages into the database, returning the first error encountered.
func (m *Model) StoreMessages(ctx context.Context, msgs []Message) error {
	for _, msg := range msgs {
		if err := m.db.Set(msg.MessageID, msg); err != nil {
			return err
		}
	}
	return nil
}

// Returns the message with the given ID and all of its replies.
func (m *Model) GetMessageAndReplies(ctx context.Context, msgID int64) ([]Message, error) {
	msgs := []Message{}
	id := msgID
	for i := 0; i < maxLinkedMessages; i++ {
		msg := Message{}
		if err := m.db.Get(id, &msg); err != nil {
			if err == pebble.ErrNotFound {
				break
			}
			return nil, err
		}
		// prepend replies to the list
		msgs = append([]Message{msg}, msgs...)
		if msg.ReplyID == nil {
			break
		}
		id = *msg.ReplyID
	}
	return msgs, nil
}
