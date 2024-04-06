package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tiagocosta/video-enconder/internal/entity"
)

func TestNewJob(t *testing.T) {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "path")

	job, err := entity.NewJob("path", "Converted", video)

	assert.NotNil(t, job)
	assert.Nil(t, err)
}
