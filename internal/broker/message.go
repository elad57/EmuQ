package broker

import (
	"time"
	"github.com/google/uuid"
)

func NewBrokerMessage(body []byte) *Message {
	return &Message{
		ID: uuid.New().String(),
		Body: body,
		Timestamp: time.Now(),
	}
}

