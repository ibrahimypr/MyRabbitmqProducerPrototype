package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection	*amqp.Connection
	Channel		*amqp.Channel
	Queues		map[string]amqp.Queue
}

func InitalizeRabbitMQ() *RabbitMQ {
	rmq, err := NewRabbitServer()
	if err != nil {
		log.Fatal(err)
	}
	return rmq
}

func NewRabbitServer() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &RabbitMQ{
		Connection: conn,
		Channel: ch,
		Queues: make(map[string]amqp.Queue),
	}, nil
}

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	if queue, exists := r.Queues[name]; exists {
		return queue, nil
	}
	q, err := r.Channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return q, err
	}
	r.Queues[name] = q
	return q, nil
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Connection.Close()
}

func (r *RabbitMQ) PrepareMessage(ctx context.Context, queueName, body string) error {
	queue, err := r.DeclareQueue(queueName)
	if err != nil {
		return err
	}
	return r.Channel.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: 		 []byte(body),
		})
}



func SendMessage(rmq *RabbitMQ, queueName, body string) {
	ctx, cancel := context.WithTimeout(context.Background(), 8 * time.Second)
	defer cancel()

	err := rmq.PrepareMessage(ctx, queueName, body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sent %s to %s\n", body, queueName)
}