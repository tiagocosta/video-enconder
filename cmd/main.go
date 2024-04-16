package main

import (
	"sync"

	"github.com/streadway/amqp"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/event/handler"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/framework/rabbitmq"
	"go.uber.org/zap"
)

var (
	eventDispatcher events.EventDispatcher
)

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()

	// Log messages
	logger.Info("This is an info message", zap.String("key", "value"))

	rabbitMQChannel := rabbitmq.OpenRabbitMQChannel()
	defer rabbitMQChannel.Close()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("JobCompleted", &handler.JobCompletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	go consumeQueue(rabbitMQChannel)
}

func consumeQueue(rabbitMQChannel *amqp.Channel) {
	msgs := make(chan amqp.Delivery)

	go rabbitmq.Consume(rabbitMQChannel, msgs, "videos")

	for msg := range msgs {
		evt := event.NewVideoRequested()
		evt.SetPayload(msg.Body)
		handler := handler.NewVideoRequestedHandler(
			&eventDispatcher,
			database.NewVideoRepository(nil),
			database.NewJobRepository(nil),
		)
		handler.Handle(evt, &sync.WaitGroup{})
		msg.Ack(false)
	}
}
