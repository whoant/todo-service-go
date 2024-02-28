package pubsub

import (
	"context"
	"log"
	"sync"

	"todo-service/common"
)

type localPubSub struct {
	name         string
	messageQueue chan *Message
	mapChannel   map[Topic][]chan *Message
	locker       *sync.RWMutex
}

func NewPubSub(name string) *localPubSub {
	pb := &localPubSub{
		name:         name,
		messageQueue: make(chan *Message, 10000),
		mapChannel:   make(map[Topic][]chan *Message),
		locker:       new(sync.RWMutex),
	}

	return pb
}

func (ps *localPubSub) Publish(ctx context.Context, topic Topic, data *Message) error {
	data.SetChannel(topic)

	go func() {
		defer common.Recovery()
		ps.messageQueue <- data
		log.Println("New message published:", data.String())
	}()

	return nil
}

func (ps *localPubSub) Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, unsubscribe func()) {
	c := make(chan *Message)
	ps.locker.Lock()

	if val, ok := ps.mapChannel[topic]; ok {
		val = append(ps.mapChannel[topic], c)
		ps.mapChannel[topic] = val
	} else {
		ps.mapChannel[topic] = []chan *Message{c}
	}

	ps.locker.Unlock()

	return c, func() {
		log.Println("Unsubscribe")
		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				if chans[i] == c {
					chans = append(chans[:i], chans[i+1:]...)
					ps.locker.Lock()
					ps.mapChannel[topic] = chans
					ps.locker.Unlock()
					close(c)
					break
				}
			}
		}

	}
}

func (ps *localPubSub) run() error {
	go func() {
		defer common.Recovery()
		for {
			msg := <-ps.messageQueue
			log.Println("Message dequeue:", msg.String())
			ps.locker.RLock()

			if subs, ok := ps.mapChannel[msg.Channel()]; ok {
				for i := range subs {
					go func(c chan *Message) {
						defer common.Recovery()
						c <- msg
					}(subs[i])
				}
			}
			ps.locker.RUnlock()
		}
	}()

	return nil
}

func (ps *localPubSub) GetPrefix() string {
	return ps.name
}

func (ps *localPubSub) Get() interface{} {
	return ps
}

func (ps *localPubSub) Name() string {
	return ps.name
}

func (*localPubSub) InitFlags() {
}

func (*localPubSub) Configure() error {
	return nil
}

func (ps *localPubSub) Run() error {
	return ps.run()
}

func (*localPubSub) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}
