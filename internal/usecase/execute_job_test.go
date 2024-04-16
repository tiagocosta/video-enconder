package usecase_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/usecase"

	_ "github.com/mattn/go-sqlite3"
)

type ExecuteJobTestSuite struct {
	suite.Suite
	DB              *sql.DB
	EventDispatcher *events.EventDispatcher
}

func (suite *ExecuteJobTestSuite) SetupSuite() {
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
	suite.EventDispatcher = events.NewEventDispatcher()
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestExecuteJobTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteJobTestSuite))
}

func (suite *ExecuteJobTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func (suite *ExecuteJobTestSuite) TestExecuteJob() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "example.mp4")
	videoRepository := database.NewVideoRepository(suite.DB)
	videoRepository.Save(video)

	job, _ := entity.NewJob("encoder_example_test", entity.Pending, video)
	jobRepository := database.NewJobRepository(suite.DB)
	jobRepository.Save(job)

	useCaseExecuteJob := usecase.NewExecuteJobUseCase(job, videoRepository, jobRepository, event.NewJobCompleted(), suite.EventDispatcher)
	inputExecuteJobDTO := usecase.ExecuteJobInputDTO{
		BucketName: job.OutputBucketPath,
		FilePath:   video.FilePath,
		VideoID:    video.ID,
	}

	err := useCaseExecuteJob.Execute(inputExecuteJobDTO)

	suite.NoError(err)
}
