package messaging_system

type RabbitMQ struct{}

func NewRabbitMq() *RabbitMQ {
	return &RabbitMQ{}
}

func (r *RabbitMQ) Publish(topic string, msg string) {}
