package rabbitmq

func NewProvider(host, port, user, pass, queue, exchange, routingKey string) (*RabbitMQ, error) {
	rabbitmq, err := initRabbitMQ(host, port, user, pass, queue)
	if err != nil {
		return nil, err
	}

	if err := rabbitmq.channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	if err := rabbitmq.channel.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
		return nil, err
	}

	rabbitmq.exchange = exchange
	rabbitmq.routingKey = routingKey

	return rabbitmq, nil
}
