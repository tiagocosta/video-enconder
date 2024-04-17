package main

import (
	"database/sql"

	"github.com/streadway/amqp"
	"github.com/tiagocosta/video-enconder/configs"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/event/handler"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/framework/rabbitmq"

	_ "github.com/go-sql-driver/mysql"
)

var (
	eventDispatcher *events.EventDispatcher
	db              *sql.DB
)

func main() {
	// // Create a logger
	// logger, _ := zap.NewProduction()

	// // Log messages
	// logger.Info("This is an info message", zap.String("key", "value"))

	configs.LoadConfig(".")

	db = database.SqlDB()
	defer db.Close()

	rabbitMQChannel := rabbitmq.OpenRabbitMQChannel()
	defer rabbitMQChannel.Close()

	eventDispatcher = events.NewEventDispatcher()
	eventDispatcher.Register("JobCompleted", &handler.JobCompletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	consumeQueue(rabbitMQChannel)
}

func consumeQueue(rabbitMQChannel *amqp.Channel) {
	msgs := make(chan amqp.Delivery)

	go rabbitmq.Consume(rabbitMQChannel, msgs, "videos")

	for msg := range msgs {
		evt := event.NewVideoRequested()
		evt.SetPayload(msg.Body)
		handler := handler.NewVideoRequestedHandler(
			eventDispatcher,
			database.NewVideoRepository(db),
			database.NewJobRepository(db),
		)
		handler.Handle(evt)
		msg.Ack(false)
	}
}
