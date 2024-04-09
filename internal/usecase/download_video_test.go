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

type DownloadVideoTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func (suite *DownloadVideoTestSuite) SetupSuite() {
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

func TestDownloadVideoTestSuite(t *testing.T) {
	suite.Run(t, new(DownloadVideoTestSuite))
}

func (suite *DownloadVideoTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func (suite *DownloadVideoTestSuite) TestVideoDownload() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")

	useCase := usecase.NewDownloadVideoUseCase()
	input := usecase.DownloadVideoInputDTO{
		BucketName: "encoder_example_test",
		FilePath:   video.FilePath,
		VideoID:    video.ID,
	}

	err := useCase.Execute(input)
	suite.NoError(err)

	useCaseCleanVideo := usecase.NewCleanVideoUseCase()
	inputCleanVideo := usecase.CleanVideoInputDTO{
		VideoID: video.ID,
	}
	useCaseCleanVideo.Execute(inputCleanVideo)
}
