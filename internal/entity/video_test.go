package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tiagocosta/video-enconder/internal/entity"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	_, err := entity.NewVideo("", "", "")
	assert.NotNil(t, err)
}

func TestValidateIfVideoIdIsNotUuid(t *testing.T) {
	_, err := entity.NewVideo(uuid.NewString(), "resourceid", "path")
	assert.Nil(t, err)
}
