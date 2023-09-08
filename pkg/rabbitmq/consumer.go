package rabbitmq

func NewConsumer(host, port, user, pass, queue string) (*RabbitMQ, error) {
	rabbitmq, err := initRabbitMQ(host, port, user, pass, queue)
	if err != nil {
		return nil, err
	}

	rabbitmq.consumerChan, err = rabbitmq.channel.Consume(queue, "", true, false, false, true, nil)
	if err != nil {
		return nil, err
	}

	return rabbitmq, nil
}
