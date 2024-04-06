package entity

import (
	"errors"
	"time"
)

type Video struct {
	ID         string
	ResourceID string
	FilePath   string
	CreatedAt  time.Time
	Jobs       []*Job
}

func NewVideo(id string, resourceID string, filePath string) (*Video, error) {
	video := &Video{
		ID:         id,
		ResourceID: resourceID,
		FilePath:   filePath,
		CreatedAt:  time.Now(),
	}

	err := video.Validate()

	if err != nil {
		return nil, err
	}

	return video, nil
}

func (video *Video) Validate() error {
	if video.ID == "" {
		return errors.New("invalid id")
	}
	if video.ResourceID == "" {
		return errors.New("invalid resource_id")
	}
	if video.FilePath == "" {
		return errors.New("invalid file path")
	}

	return nil
}
