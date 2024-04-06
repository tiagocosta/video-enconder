package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tiagocosta/video-enconder/internal/entity"
)

func TestNewJob(t *testing.T) {
	video := entity.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := entity.NewJob("path", "Converted", video)

	assert.NotNil(t, job)
	assert.Nil(t, err)
}
