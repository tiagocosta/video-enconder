package usecase

import (
	"github.com/google/uuid"
	"github.com/tiagocosta/video-enconder/configs"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/pkg/encoder"
)

type ExecuteJobInputDTO struct {
	ResourceID string `json:"resource_id"`
	FilePath   string `json:"file_path"`
}

type ExecuteJobUseCase struct {
	VideoRepository *database.VideoRepository
	JobRepository   *database.JobRepository
	JobCompleted    *event.JobCompleted
	EventDispatcher *events.EventDispatcher
	Encoder         encoder.VideoEncoder
}

func NewExecuteJobUseCase(
	videoRepository *database.VideoRepository,
	jobRepository *database.JobRepository,
	jobCompleted *event.JobCompleted,
	eventDispatcher *events.EventDispatcher,
	encoder encoder.VideoEncoder,
) *ExecuteJobUseCase {
	return &ExecuteJobUseCase{
		VideoRepository: videoRepository,
		JobRepository:   jobRepository,
		JobCompleted:    jobCompleted,
		EventDispatcher: eventDispatcher,
		Encoder:         encoder,
	}
}

func (uc *ExecuteJobUseCase) Execute(input ExecuteJobInputDTO) error {
	video, _ := entity.NewVideo(uuid.NewString(), input.ResourceID, input.FilePath)
	uc.VideoRepository.Save(video)

	bucketName := configs.Config().BucketName
	job, _ := entity.NewJob(bucketName, entity.Pending, video)
	uc.JobRepository.Save(job)

	job.StartVideoDownload()
	err := uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}
	err = uc.Encoder.Download(job.Video.FilePath, job.VideoID)
	if err != nil {
		return uc.failJob(job, err)
	}

	job.StartVideoFragmentation()
	err = uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}
	err = uc.Encoder.Fragment(job.VideoID)
	if err != nil {
		return uc.failJob(job, err)
	}

	job.StartVideoEncoding()
	err = uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}
	err = uc.Encoder.Encode(job.VideoID)
	if err != nil {
		return uc.failJob(job, err)
	}

	job.StartVideoUpload()
	err = uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}
	err = uc.Encoder.Upload(job.VideoID)
	if err != nil {
		return uc.failJob(job, err)
	}

	job.CleanupVideoFiles()
	err = uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}
	err = uc.Encoder.CleanupFiles(job.VideoID)
	if err != nil {
		return uc.failJob(job, err)
	}

	job.Complete()
	err = uc.updateJob(job)
	if err != nil {
		return uc.failJob(job, err)
	}

	return nil
}

func (uc *ExecuteJobUseCase) updateJob(job *entity.Job) error {
	err := uc.JobRepository.Update(job)

	if err != nil {
		return err
	}

	return nil
}

func (uc *ExecuteJobUseCase) failJob(job *entity.Job, error error) error {
	job.Fail()
	err := uc.updateJob(job)

	if err != nil {
		return err
	}

	return error
}
