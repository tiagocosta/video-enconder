package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
)

type JobCompletedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewJobCompletedHandler(rabbitMQChannel *amqp.Channel) *JobCompletedHandler {
	return &JobCompletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *JobCompletedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Job completed: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"jobs",       // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish
	)
}
