package messaging_system

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn     *amqp.Connection
	Url      string
	Protocol string
	Username string
	Password string
}

func (r *RabbitMQ) Connect() error {
	connString := r.Protocol + "://" + r.Username + ":" + r.Password + "@" + r.Url
	conn, err := amqp.Dial(connString)
	if err != nil {
		return err
	}
	r.conn = conn
	return r.init()
}

func (r *RabbitMQ) init() error {
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	err = channel.ExchangeDeclare(
		os.Getenv("RABBITMQ_EXCHANGE_NAME"),
		os.Getenv("RABBITMQ_EXCHANGE_TYPE"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE_NAME"),
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = channel.QueueBind(
		queue.Name,
		os.Getenv("RABBITMQ_ROUTING_KEY"),
		os.Getenv("RABBITMQ_EXCHANGE_NAME"),
		false,
		nil,
	)
	return err
}

func (r *RabbitMQ) PublishWithCtx(msg []byte) error {
	ch, err := r.conn.Channel()
	if err != nil {

	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		os.Getenv("RABBITMQ_EXCHANGE_NAME"),
		os.Getenv("RABBITMQ_QUEUE_NAME"),
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         msg,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
