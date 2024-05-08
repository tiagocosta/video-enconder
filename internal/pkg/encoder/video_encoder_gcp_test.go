package encoder_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/tiagocosta/video-enconder/configs"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/pkg/encoder"

	_ "github.com/mattn/go-sqlite3"
)

type VideoUtilsTestSuite struct {
	suite.Suite
	Encoder *encoder.VideoEncoderGCP
}

func (suite *VideoUtilsTestSuite) SetupSuite() {
	configs.LoadConfig("../../../")
	suite.Encoder = &encoder.VideoEncoderGCP{}
}

func TestVideoUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(VideoUtilsTestSuite))
}

func (suite *VideoUtilsTestSuite) TestDownload() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")

	err := suite.Encoder.Download(video.FilePath, video.ID)
	suite.NoError(err)

	suite.Encoder.CleanupFiles(video.ID)
}

func (suite *VideoUtilsTestSuite) TestFragment() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")

	err := suite.Encoder.Download(video.FilePath, video.ID)
	suite.NoError(err)

	err = suite.Encoder.Fragment(video.ID)
	suite.NoError(err)

	suite.Encoder.CleanupFiles(video.ID)
}

func (suite *VideoUtilsTestSuite) TestEncode() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")

	err := suite.Encoder.Download(video.FilePath, video.ID)
	suite.NoError(err)

	err = suite.Encoder.Fragment(video.ID)
	suite.NoError(err)

	err = suite.Encoder.Encode(video.ID)
	suite.NoError(err)

	suite.Encoder.CleanupFiles(video.ID)
}

func (suite *VideoUtilsTestSuite) TestUpload() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")

	err := suite.Encoder.Download(video.FilePath, video.ID)
	suite.NoError(err)

	err = suite.Encoder.Fragment(video.ID)
	suite.NoError(err)

	err = suite.Encoder.Encode(video.ID)
	suite.NoError(err)

	err = suite.Encoder.Upload(video.ID)
	suite.NoError(err)

	err = suite.Encoder.CleanupFiles(video.ID)
	suite.NoError(err)
}
