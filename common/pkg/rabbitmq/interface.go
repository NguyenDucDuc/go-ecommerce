package rabbitmq

type IRabbitMQService interface {
	Publish(exchange, routingKey string, body interface{}) error
	Consume(queueName, routingKey, exchange string, handler func([]byte))
	Close()
}