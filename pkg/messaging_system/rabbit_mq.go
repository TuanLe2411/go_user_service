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

	args := amqp.Table{
		"x-message-ttl": int32(5000),
		"x-max-length":  int32(500),
	}
	queue, err := channel.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE_NAME"),
		false,
		false,
		true,
		false,
		args,
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

func (r *RabbitMQ) Publish(msg []byte) error {
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
			Expiration:   "5000",
		},
	)
	if err != nil {
		log.Println("Error when publish message: " + string(msg) + " ,err: " + err.Error())
		return err
	}
	log.Println("Publish message successfully, message: " + string(msg))
	return nil
}
