package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/usecase"
)

type VideoRequestedInputDTO struct {
	ResourceID string
	FilePath   string
}

type VideoRequestedHandler struct {
	EventDispatcher *events.EventDispatcher
	VideoRepository *database.VideoRepository
	JobRepository   *database.JobRepository
}

func NewVideoRequestedHandler(
	eventDispatcher *events.EventDispatcher,
	videoRepository *database.VideoRepository,
	jobRepository *database.JobRepository,
) *VideoRequestedHandler {
	return &VideoRequestedHandler{
		EventDispatcher: eventDispatcher,
		VideoRepository: videoRepository,
		JobRepository:   jobRepository,
	}
}

func (h *VideoRequestedHandler) Handle(evt events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("video requested: %v", evt.GetPayload())
	inputDTO := VideoRequestedInputDTO{}
	json.Unmarshal(evt.GetPayload().([]byte), &inputDTO)

	video, _ := entity.NewVideo(uuid.NewString(), inputDTO.ResourceID, inputDTO.FilePath)
	h.VideoRepository.Save(video)

	bucketName := os.Getenv("BUCKET_NAME")
	job, _ := entity.NewJob(bucketName, entity.Pending, video)
	h.JobRepository.Save(job)

	jobCompleted := event.NewJobCompleted()

	useCaseExecuteJob := usecase.NewExecuteJobUseCase(job, h.VideoRepository, h.JobRepository, jobCompleted, h.EventDispatcher)
	inputExecuteJobDTO := usecase.ExecuteJobInputDTO{
		BucketName: job.OutputBucketPath,
		FilePath:   video.FilePath,
		VideoID:    video.ID,
	}

	useCaseExecuteJob.Execute(inputExecuteJobDTO)
}
