package consumer

import (
	"github.com/streadway/amqp"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/event/handler"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/framework/rabbitmq"
	"github.com/tiagocosta/video-enconder/internal/pkg/encoder"
)

type VideoConsumer struct {
	Channel         *amqp.Channel
	EventDispatcher *events.EventDispatcher
	VideoRepository *database.VideoRepository
	JobRepository   *database.JobRepository
	Encoder         encoder.VideoEncoder
}

func (consumer *VideoConsumer) Initialize() {
	rabbitmq.DeclareExchange(consumer.Channel, "dlx", "direct", nil)

	queueArgs := make(amqp.Table)
	queueArgs["x-dead-letter-exchange"] = "dlx"

	rabbitmq.DeclareQueue(consumer.Channel, "videos", queueArgs)
	rabbitmq.DeclareQueue(consumer.Channel, "videos-failed", nil)
	rabbitmq.DeclareQueue(consumer.Channel, "videos-result", nil)

	rabbitmq.BindQueue(consumer.Channel, "dlx", "videos-failed", "jobs")
	rabbitmq.BindQueue(consumer.Channel, "amq.direct", "videos-result", "jobs")
}

func (consumer *VideoConsumer) ConsumeQueue() {
	msgs := make(chan amqp.Delivery)

	go rabbitmq.Consume(consumer.Channel, msgs, "videos")

	for msg := range msgs {
		evt := event.NewVideoRequested()
		evt.SetPayload(msg.Body)
		handler := handler.NewVideoRequestedHandler(
			consumer.EventDispatcher,
			consumer.VideoRepository,
			consumer.JobRepository,
			consumer.Encoder,
		)
		err := handler.Handle(evt)
		if err != nil {
			msg.Nack(false, false)
		} else {
			msg.Ack(false)
		}
	}
}
