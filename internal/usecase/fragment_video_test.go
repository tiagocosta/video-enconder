package usecase_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/usecase"

	_ "github.com/mattn/go-sqlite3"
)

type FragmentVideoTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func (suite *FragmentVideoTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec(
		`CREATE TABLE video (
			id varchar(255) NOT NULL,
			resource_id varchar(255) NOT NULL,
			file_path varchar(255) NOT NULL,
			created_at datetime NOT NULL,
			PRIMARY KEY (id)
			)`)
	db.Exec(
		`CREATE TABLE job (
				id varchar(255) NOT NULL,
				output_bucket_path varchar(255) NOT NULL,
				status varchar(255) NOT NULL,
				video_id varchar(255) NOT NULL,
				created_at datetime NOT NULL,
				updated_at datetime,
				PRIMARY KEY (id)
		)`)
	suite.DB = db
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestFragmentVideoTestSuite(t *testing.T) {
	suite.Run(t, new(FragmentVideoTestSuite))
}

func (suite *FragmentVideoTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func (suite *FragmentVideoTestSuite) TestVideoFragment() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")
	videoRepo := database.NewVideoRepository(suite.DB)

	useCaseDownloadVideo := usecase.NewDownloadVideoUseCase(videoRepo)
	inputDownloadVideo := usecase.DownloadVideoInputDTO{
		BucketName: "encoder_example_test",
		FilePath:   video.FilePath,
		VideoID:    video.ID,
	}
	useCaseDownloadVideo.Execute(inputDownloadVideo)

	useCaseFragmentVideo := usecase.NewFragmentVideoUseCase()
	inputFragmentVideo := usecase.FragmentVideoInputDTO{
		VideoID: video.ID,
	}
	err := useCaseFragmentVideo.Execute(inputFragmentVideo)
	suite.NoError(err)
}
