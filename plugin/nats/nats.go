package nats

import (
	"context"
	"encoding/json"
	"flag"
	"time"

	"github.com/nats-io/nats.go"
	"todo-service/go-sdk/logger"
	"todo-service/pubsub"
)

type appNats struct {
	name       string
	connection *nats.Conn
	logger     logger.Logger
	url        string
}

func NewNATS(name string) *appNats {
	return &appNats{
		name: name,
	}
}

func (app *appNats) Name() string {
	return app.name
}

func (app *appNats) InitFlags() {
	flag.StringVar(&app.url, app.name+"-url", nats.DefaultURL, "URL of NATS service")
}

func (app *appNats) Configure() error {
	app.logger = logger.GetCurrent().GetLogger(app.name)
	app.logger.Infoln("Connecting to NATS service...")
	conn, err := nats.Connect(app.url, app.setupConnOptions([]nats.Option{})...)
	if err != nil {
		app.logger.Fatalln(err)
	}
	app.logger.Infoln("Connected to NATS service")
	app.connection = conn

	return nil
}

func (app *appNats) setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
		app.logger.Infof("Disconnected due to: %s, will attempt reconnect for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(conn *nats.Conn) {
		app.logger.Infof("Reconnected [%s]", conn.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(conn *nats.Conn) {
		app.logger.Infoln("Exiting: %v", conn.LastError())
	}))
	return opts
}

func (app *appNats) Run() error {
	return app.Configure()
}

func (app *appNats) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (app *appNats) GetPrefix() string {
	return app.name
}

func (app *appNats) Get() interface{} {
	return app
}

func (app *appNats) Publish(ctx context.Context, channel pubsub.Topic, data *pubsub.Message) error {
	msgData, err := json.Marshal(data.Data())
	if err != nil {
		app.logger.Errorln(err)
		return err
	}

	if err := app.connection.Publish(string(channel), msgData); err != nil {
		app.logger.Errorln(err)
		return err
	}

	return nil
}

func (app *appNats) Subscribe(ctx context.Context, channel pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	msgChan := make(chan *pubsub.Message)

	sub, err := app.connection.Subscribe(string(channel), func(msg *nats.Msg) {
		msgData := make(map[string]interface{})
		_ = json.Unmarshal(msg.Data, &msgData)
		appMsg := pubsub.NewMessage(msgData)
		appMsg.SetChannel(channel)
		appMsg.SetAckFunc(func() error {
			return msg.Ack()
		})

		msgChan <- appMsg
	})

	if err != nil {
		app.logger.Errorln(err)
	}

	return msgChan, func() {
		_ = sub.Unsubscribe()
	}
}
