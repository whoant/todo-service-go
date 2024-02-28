package pubsub

import (
	"context"
	"fmt"
	"time"
)

type Topic string

type Message struct {
	id        string
	channel   Topic
	data      interface{}
	createdAt time.Time
	ackFunc   func() error
}

type PubSub interface {
	Publish(ctx context.Context, channel Topic, data *Message) error
	Subscribe(ctx context.Context, channel Topic) (ch <-chan *Message, close func())
}

func (m *Message) SetAckFunc(f func() error) {
	m.ackFunc = f
}

func (m *Message) Ack() error {
	return m.ackFunc()
}

func NewMessage(data interface{}) *Message {
	now := time.Now().UTC()

	return &Message{
		id:        fmt.Sprintf("%d", now.UnixNano()),
		data:      data,
		createdAt: now,
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("Message %s value %v", m.channel, m.data)
}

func (m *Message) Channel() Topic {
	return m.channel
}

func (m *Message) SetChannel(channel Topic) {
	m.channel = channel
}

func (m *Message) Data() interface{} {
	return m.data
}
