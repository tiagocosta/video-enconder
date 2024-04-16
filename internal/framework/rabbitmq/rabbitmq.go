package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func OpenRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func Consume(rabbitMQChannel *amqp.Channel, out chan<- amqp.Delivery, queue string) error {
	rabbitMQQueue, err := rabbitMQChannel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "failed to declare queue", err)
	}

	msgs, err := rabbitMQChannel.Consume(
		rabbitMQQueue.Name,
		"go-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for msg := range msgs {
		out <- msg
	}
	return nil
}
