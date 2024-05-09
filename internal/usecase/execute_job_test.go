package usecase_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tiagocosta/video-enconder/configs"
	"github.com/tiagocosta/video-enconder/internal/event"
	"github.com/tiagocosta/video-enconder/internal/framework/database"
	"github.com/tiagocosta/video-enconder/internal/framework/events"
	"github.com/tiagocosta/video-enconder/internal/pkg/encoder"
	"github.com/tiagocosta/video-enconder/internal/usecase"

	_ "github.com/mattn/go-sqlite3"
)

type ExecuteJobTestSuite struct {
	suite.Suite
	DB              *sql.DB
	EventDispatcher *events.EventDispatcher
	Encoder         encoder.VideoEncoder
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
	suite.Encoder = &encoder.VideoEncoderGCP{}
	configs.LoadConfig("../../")
}

func TestExecuteJobTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteJobTestSuite))
}

func (suite *ExecuteJobTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func (suite *ExecuteJobTestSuite) TestExecuteJob() {
	videoRepository := database.NewVideoRepository(suite.DB)
	jobRepository := database.NewJobRepository(suite.DB)

	useCaseExecuteJob := usecase.NewExecuteJobUseCase(
		videoRepository,
		jobRepository,
		event.NewJobCompleted(),
		event.NewJobError(),
		suite.EventDispatcher,
		suite.Encoder,
	)
	inputExecuteJobDTO := usecase.ExecuteJobInputDTO{
		FilePath:   "example.mp4",
		ResourceID: "resource_id",
	}

	go useCaseExecuteJob.Execute(inputExecuteJobDTO)

	// suite.NoError(err)
}
