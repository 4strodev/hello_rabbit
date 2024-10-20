package infrastructure

import (
	"log/slog"
	"os"

	"github.com/4strodev/wiring/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitComponent struct {
	connection *amqp.Connection
	Logger     *slog.Logger
}

// Init implements components.Component.
func (r *RabbitComponent) Init(container pkg.Container) error {
	var err error
	err = container.Fill(r)
	if err != nil {
		return err
	}

	conn, err := amqp.Dial(os.Getenv("RABBIT_URL"))
	if err != nil {
		return err
	}
	r.connection = conn

	err = container.Singleton(func() *amqp.Connection {
		return conn
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitComponent) OnShutdown() error {
	r.Logger.Info("")
	return r.connection.Close()
}
