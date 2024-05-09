package main

import (
	"github.com/tiagocosta/video-enconder/configs"
	"github.com/tiagocosta/video-enconder/internal/event/consumer"
	"github.com/tiagocosta/video-enconder/internal/event/handler"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/framework/rabbitmq"
	"github.com/tiagocosta/video-enconder/internal/pkg/encoder"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs.LoadConfig(".")

	db := database.SqlDB()
	defer db.Close()

	rabbitMQChannel := rabbitmq.OpenRabbitMQChannel()
	defer rabbitMQChannel.Close()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("JobCompleted", &handler.JobCompletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("JobError", &handler.JobErrorHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	videoConsumer := consumer.VideoConsumer{
		Channel:         rabbitMQChannel,
		EventDispatcher: eventDispatcher,
		VideoRepository: database.NewVideoRepository(db),
		JobRepository:   database.NewJobRepository(db),
		Encoder:         &encoder.VideoEncoderGCP{},
	}
	videoConsumer.Initialize()
	videoConsumer.ConsumeQueue()
}
