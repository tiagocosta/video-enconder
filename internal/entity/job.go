package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID               string
	OutputBucketPath string
	Status           string
	Video            *Video
	VideoID          string
	Error            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewJob(output string, status string, video *Video) (*Job, error) {
	job := Job{
		ID:               uuid.NewString(),
		OutputBucketPath: output,
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
