package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	conn *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQ khởi tạo kết nối và trả về một instance
func NewRabbitMQ(url string) (IRabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Khai báo Exchange chung cho hệ thống (ví dụ: Topic Exchange)
	err = ch.ExchangeDeclare(
		"topic_exchange", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, err
	}
	log.Println("✅ kết nối thành công tới RabbitMQ")

	return &RabbitMQService{
		conn:    conn,
		channel: ch,
	}, nil
}

// Publish gửi tin nhắn vào một routing key cụ thể
func (r *RabbitMQService) Publish(exchange, routingKey string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(context.Background(),
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

func (r *RabbitMQService) Consume(queueName, routingKey, exchange string, handler func([]byte)) {
	// 1. Khai báo Queue
	q, _ := r.channel.QueueDeclare(queueName, true, false, false, false, nil)

	// 2. Bind Queue với Exchange
	r.channel.QueueBind(q.Name, routingKey, exchange, false, nil)

	// 3. Lắng nghe
	msgs, _ := r.channel.Consume(q.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			handler(d.Body) // Gọi callback xử lý logic
		}
	}()
}

// Close để đóng kết nối khi tắt app
func (r *RabbitMQService) Close() {
	r.channel.Close()
	r.conn.Close()
}