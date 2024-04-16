package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Pending     Status = "PENDING"
	Starting    Status = "STARTING"
	Downloading Status = "DOWNLOADING"
	Fragmenting Status = "FRAGMENTING"
	Encoding    Status = "ENCODING"
	Uploading   Status = "UPLOADING"
	Finishing   Status = "FINISHING"
	Completed   Status = "COMPLETED"
	Failed      Status = "FAILED"
)

type Job struct {
	ID               string
	OutputBucketPath string
	Status           Status
	Video            *Video
	VideoID          string
	Error            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewJob(outputBucket string, status Status, video *Video) (*Job, error) {
	job := Job{
		ID:               uuid.NewString(),
		OutputBucketPath: outputBucket,
		Status:           status,
		Video:            video,
		VideoID:          video.ID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (job *Job) Validate() error {
	if job.ID == "" {
		return errors.New("invalid id")
	}
	err := uuid.Validate(job.ID)
	if err != nil {
		return errors.New("invalid uuid format")
	}
	if job.VideoID == "" {
		return errors.New("invalid video_id")
	}
	if job.OutputBucketPath == "" {
		return errors.New("invalid output bucket path")
	}
	if job.Status == "" {
		return errors.New("invalid status")
	}

	return nil
}

func (job *Job) StartVideoDownload() {
	job.Status = Downloading
}

func (job *Job) StartVideoFragmentation() {
	job.Status = Fragmenting
}

func (job *Job) StartVideoEncoding() {
	job.Status = Encoding
}

func (job *Job) StartVideoUpload() {
	job.Status = Uploading
}

func (job *Job) CleanupVideoFiles() {
	job.Status = Finishing
}

func (job *Job) Complete() {
	job.Status = Completed
}

func (job *Job) Fail() {
	job.Status = Failed
}
