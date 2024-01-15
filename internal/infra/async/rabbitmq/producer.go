package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type producer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (p producer) declareExchange(ctx context.Context, name string) error {
	var (
		kind                  = "topic"
		durable               = true
		autoDelete            = false
		internal              = false
		noWait                = false
		args       amqp.Table = nil
	)

	return p.ch.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (p producer) ensureExchanges(ctx context.Context) error {
	if err := p.ensureStatisticExchange(ctx); err != nil {
		return fmt.Errorf("Statistic Exchange: %w", err)
	}

	return nil
}

func (p producer) getChannel() (*amqp.Channel, error) {
	if p.ch != nil {
		if !p.ch.IsClosed() {
			return p.ch, nil
		}

		_ = p.ch.Close()
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Get Channel: %w", err)
	}

	p.ch = ch
	return p.ch, nil
}

func (p producer) Close() {
	defer p.conn.Close()
	defer p.ch.Close()
}

func NewProducer(ctx context.Context, rabbitmqURI string) (*producer, error) {
	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		defer conn.Close()
		return nil, err
	}

	p := &producer{
		conn: conn,
		ch:   ch,
	}
	return p, p.ensureExchanges(ctx)
}
