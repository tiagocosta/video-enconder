package usecase_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/usecase"

	_ "github.com/mattn/go-sqlite3"
)

type CleanVideoTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func (suite *CleanVideoTestSuite) SetupSuite() {
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

func TestCleanVideoTestSuite(t *testing.T) {
	suite.Run(t, new(CleanVideoTestSuite))
}

func (suite *CleanVideoTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func (suite *CleanVideoTestSuite) TestVideoClean() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")

	useCaseDownloadVideo := usecase.NewDownloadVideoUseCase()
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
	useCaseFragmentVideo.Execute(inputFragmentVideo)

	useCaseEncodeVideo := usecase.NewEncodeVideoUseCase()
	inputEncodeVideo := usecase.EncodeVideoInputDTO{
		VideoID: video.ID,
	}
	useCaseEncodeVideo.Execute(inputEncodeVideo)

	useCaseCleanVideo := usecase.NewCleanVideoUseCase()
	inputCleanVideo := usecase.CleanVideoInputDTO{
		VideoID: video.ID,
	}
	err := useCaseCleanVideo.Execute(inputCleanVideo)
	suite.NoError(err)
}
