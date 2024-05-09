package handler

import (
	"encoding/json"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/framework/rabbitmq"
)

type JobErrorHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewJobErrorHandler(rabbitMQChannel *amqp.Channel) *JobErrorHandler {
	return &JobErrorHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *JobErrorHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) error {
	defer wg.Done()

	jsonOutput, _ := json.Marshal(event.GetPayload())
	rabbitmq.Publish(h.RabbitMQChannel, jsonOutput, "dlx", "jobs")

	return nil
}
