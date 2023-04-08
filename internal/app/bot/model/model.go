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

// Returns all replies of the given message.
func (m *Model) GetReplies(ctx context.Context, msgID int64) ([]Message, error) {
	msgs := []Message{}
	id := msgID
	for i := 0; i < maxLinkedMessages; i++ {
		msg := Message{}
		err := m.db.Get(id, &msg)
		if err == pebble.ErrNotFound {
			break
		}
		if err != nil {
			return nil, err
		}
		// prepend replies to the list
		if i != 0 {
			msgs = append([]Message{msg}, msgs...)
		}
		if msg.ReplyID == nil {
			break
		}
		id = *msg.ReplyID
	}
	return msgs, nil
}
