package handler

import (
	"encoding/json"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/framework/rabbitmq"
)

type JobCompletedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewJobCompletedHandler(rabbitMQChannel *amqp.Channel) *JobCompletedHandler {
	return &JobCompletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *JobCompletedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) error {
	defer wg.Done()

	jsonOutput, _ := json.Marshal(event.GetPayload())
	rabbitmq.Publish(h.RabbitMQChannel, jsonOutput, "amq.direct", "jobs")

	return nil
}
