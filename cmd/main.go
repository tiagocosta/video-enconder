package main

import (
	"database/sql"

	"github.com/tiagocosta/video-enconder/configs"
	"github.com/tiagocosta/video-enconder/internal/event/consumer"
	"github.com/tiagocosta/video-enconder/internal/event/handler"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/framework/rabbitmq"
	"github.com/tiagocosta/video-enconder/internal/pkg/encoder"

	_ "github.com/go-sql-driver/mysql"
)

var (
	eventDispatcher *events.EventDispatcher
	db              *sql.DB
	videoEncoder    encoder.VideoEncoder
)

func main() {
	configs.LoadConfig(".")

	db = database.SqlDB()
	defer db.Close()

	rabbitMQChannel := rabbitmq.OpenRabbitMQChannel()
	defer rabbitMQChannel.Close()

	eventDispatcher = events.NewEventDispatcher()
	eventDispatcher.Register("JobCompleted", &handler.JobCompletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	videoConsumer := consumer.VideoConsumer{
		Channel:         rabbitMQChannel,
		EventDispatcher: eventDispatcher,
		VideoRepository: database.NewVideoRepository(db),
		JobRepository:   database.NewJobRepository(db),
		Encoder:         &encoder.VideoEncoderGCP{},
	}

	videoConsumer.ConsumeQueue()
}

// func consumeQueue(rabbitMQChannel *amqp.Channel) {
// 	msgs := make(chan amqp.Delivery)

// 	go rabbitmq.Consume(rabbitMQChannel, msgs, "videos")

// 	for msg := range msgs {
// 		evt := event.NewVideoRequested()
// 		evt.SetPayload(msg.Body)
// 		handler := handler.NewVideoRequestedHandler(
// 			eventDispatcher,
// 			database.NewVideoRepository(db),
// 			database.NewJobRepository(db),
// 			videoEncoder,
// 		)
// 		handler.Handle(evt)
// 		msg.Ack(false)
// 	}
// }
