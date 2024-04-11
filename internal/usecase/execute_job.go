package usecase

import (
	"errors"
	"os"
	"strconv"

	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
)

type ExecuteJobInputDTO struct {
	BucketName string
	FilePath   string
	VideoID    string
}

type ExecuteJobUseCase struct {
	Job           *entity.Job
	JobRepository *database.JobRepository
}

func NewExecuteJobUseCase(job *entity.Job, jobRepository *database.JobRepository) *ExecuteJobUseCase {
	return &ExecuteJobUseCase{
		Job:           job,
		JobRepository: jobRepository,
	}
}

func (uc *ExecuteJobUseCase) Execute(input ExecuteJobInputDTO) error {
	err := uc.download(input)
	if err != nil {
		return err
	}

	err = uc.fragment(input)
	if err != nil {
		return err
	}

	err = uc.encode(input)
	if err != nil {
		return err
	}

	err = uc.upload(input)
	if err != nil {
		return err
	}

	err = uc.cleanupFiles(input)
	if err != nil {
		return err
	}

	err = uc.complete()
	if err != nil {
		return err
	}

	return nil
}

func (uc *ExecuteJobUseCase) updateJob() error {
	err := uc.JobRepository.Update(uc.Job)

	if err != nil {
		return err
	}

	return nil
}

func (uc *ExecuteJobUseCase) failJob(error error) error {
	uc.Job.Fail()
	err := uc.updateJob()

	if err != nil {
		return err
	}

	return error
}

func (uc *ExecuteJobUseCase) download(input ExecuteJobInputDTO) error {
	uc.Job.StartVideoDownload()
	err := uc.updateJob()
	if err != nil {
		return uc.failJob(err)
	}

	downloadInputDTO := DownloadVideoInputDTO(input)
	err = NewDownloadVideoUseCase().Execute(downloadInputDTO)
	if err != nil {
		return uc.failJob(err)
	}

	return nil
}

func (uc *ExecuteJobUseCase) fragment(input ExecuteJobInputDTO) error {
	uc.Job.StartVideoFragmentation()
	err := uc.updateJob()
	if err != nil {
		return uc.failJob(err)
	}

	fragmentInputDTO := FragmentVideoInputDTO{
		VideoID: input.VideoID,
	}
	err = NewFragmentVideoUseCase().Execute(fragmentInputDTO)
	if err != nil {
		return uc.failJob(err)
	}

	return nil
}

func (uc *ExecuteJobUseCase) encode(input ExecuteJobInputDTO) error {
	uc.Job.StartVideoEncoding()
	err := uc.updateJob()
	if err != nil {
		return uc.failJob(err)
	}

	encodeInputDTO := EncodeVideoInputDTO{
		VideoID: input.VideoID,
	}
	err = NewEncodeVideoUseCase().Execute(encodeInputDTO)
	if err != nil {
		return uc.failJob(err)
	}

	return nil
}

func (uc *ExecuteJobUseCase) upload(input ExecuteJobInputDTO) error {
	uc.Job.StartVideoUpload()
	err := uc.updateJob()
	if err != nil {
		return uc.failJob(err)
	}

	doneUpload := make(chan string)
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	uploadInputDTO := UploadVideoInputDTO{
		BucketName:  input.BucketName,
		VideoPath:   os.Getenv("LOCAL_STORAGE_PATH") + "/" + input.VideoID,
		Concurrency: concurrency,
		DoneUpload:  doneUpload,
	}

	uploadVideoUseCase, _ := NewUploadVideoUseCase()

	go uploadVideoUseCase.Execute(uploadInputDTO)

	result := <-doneUpload
	if result != "upload completed" {
		return uc.failJob(errors.New(result))
	}

	return nil
}

func (uc *ExecuteJobUseCase) cleanupFiles(input ExecuteJobInputDTO) error {
	uc.Job.CleanupVideoFiles()
	err := uc.updateJob()
	if err != nil {
		return uc.failJob(err)
	}

	cleanInputDTO := CleanVideoInputDTO{
		VideoID: input.VideoID,
	}
	err = NewCleanVideoUseCase().Execute(cleanInputDTO)
	if err != nil {
		return uc.failJob(err)
	}

	return nil
}

func (uc *ExecuteJobUseCase) complete() error {
	uc.Job.Complete()
	err := uc.updateJob()
	if err != nil {
		return uc.failJob(err)
	}

	return nil
}