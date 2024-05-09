package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func OpenRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@172.26.0.1:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func DeclareExchange(ch *amqp.Channel, n, t string, args amqp.Table) {
	err := ch.ExchangeDeclare(
		n,
		t,
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare exchange "+n)
}

func DeclareQueue(ch *amqp.Channel, n string, args amqp.Table) {
	_, err := ch.QueueDeclare(
		n,
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare queue "+n)
}

func BindQueue(ch *amqp.Channel, e, q, k string) {
	err := ch.QueueBind(
		q,
		k,
		e,
		false,
		nil)
	failOnError(err, "Failed to bind "+e+" with "+q)
}

func Publish(ch *amqp.Channel, m []byte, e, k string) {
	err := ch.Publish(
		e,
		k,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        m,
		})
	failOnError(err, "Failed to publish on exchange "+e)
}

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err.Error())
	}
}

func Consume(rabbitMQChannel *amqp.Channel, out chan<- amqp.Delivery, queue string) error {
	// rabbitMQQueue, err := rabbitMQChannel.QueueDeclare(
	// 	queue, // name
	// 	true,  // durable
	// 	false, // delete when usused
	// 	false, // exclusive
	// 	false, // no-wait
	// 	nil,   // arguments
	// )

	// if err != nil {
	// 	log.Fatalf("%s: %s", "failed to declare queue", err)
	// }

	msgs, err := rabbitMQChannel.Consume(
		queue,
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
