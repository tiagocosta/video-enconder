package handler

import (
	"encoding/json"
	"fmt"

	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/pkg/encoder"
	"github.com/tiagocosta/video-enconder/internal/usecase"
)

type VideoRequestedInputDTO struct {
	ResourceID string `json:"resource_id"`
	FilePath   string `json:"file_path"`
}

type VideoRequestedHandler struct {
	EventDispatcher *events.EventDispatcher
	VideoRepository *database.VideoRepository
	JobRepository   *database.JobRepository
	Encoder         encoder.VideoEncoder
}

func NewVideoRequestedHandler(
	eventDispatcher *events.EventDispatcher,
	videoRepository *database.VideoRepository,
	jobRepository *database.JobRepository,
	encoder encoder.VideoEncoder,
) *VideoRequestedHandler {
	return &VideoRequestedHandler{
		EventDispatcher: eventDispatcher,
		VideoRepository: videoRepository,
		JobRepository:   jobRepository,
		Encoder:         encoder,
	}
}

func (h *VideoRequestedHandler) Handle(evt events.EventInterface) {
	fmt.Printf("video requested")
	inputDTO := VideoRequestedInputDTO{}
	json.Unmarshal(evt.GetPayload().([]byte), &inputDTO)

	jobCompleted := event.NewJobCompleted()

	useCaseExecuteJob := usecase.NewExecuteJobUseCase(
		h.VideoRepository,
		h.JobRepository,
		jobCompleted,
		h.EventDispatcher,
		h.Encoder,
	)
	inputExecuteJobDTO := usecase.ExecuteJobInputDTO{
		FilePath:   inputDTO.FilePath,
		ResourceID: inputDTO.ResourceID,
	}

	err := useCaseExecuteJob.Execute(inputExecuteJobDTO)
	if err != nil {
		fmt.Println(err)
	}
}
