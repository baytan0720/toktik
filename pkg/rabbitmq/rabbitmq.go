package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchange     string
	routingKey   string
	consumerChan <-chan amqp.Delivery
}

func initRabbitMQ(host, port, user, pass, queue string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(generateUrl(host, port, user, pass))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) Publish(msg string) error {
	return r.channel.Publish(r.exchange, r.routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}

func (r *RabbitMQ) Consume() string {
	return string((<-r.consumerChan).Body)
}

func (r *RabbitMQ) Close() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func generateUrl(host, port, user, pass string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
}
