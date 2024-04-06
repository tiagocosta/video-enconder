package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tiagocosta/video-enconder/internal/entity"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := entity.NewVideo()
	err := video.Validate()
	assert.NotNil(t, err)
}

func TestValidateIfVideoIdIsNotUuid(t *testing.T) {
	video := entity.NewVideo()
	video.ID = uuid.NewString()
	video.ResourceID = "resourceid"
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	err := video.Validate()
	assert.Nil(t, err)
}
