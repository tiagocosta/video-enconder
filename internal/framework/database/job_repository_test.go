package database_test

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/framework/database"

	_ "github.com/mattn/go-sqlite3"
)

type JobRepositoryTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func (suite *JobRepositoryTestSuite) SetupSuite() {
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
}

func TestJobRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(JobRepositoryTestSuite))
}

func (suite *JobRepositoryTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func (suite *JobRepositoryTestSuite) TestGivenAnJob_WhenSave_ThenShouldSaveJob() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "file_path")
	videoRepo := database.NewVideoRepository(suite.DB)
	videoRepo.Save(video)

	job, err := entity.NewJob("output_path", "Pending", video)
	suite.NoError(err)
	suite.NotNil(job)
	jobRepo := database.NewJobRepository(suite.DB)
	jobRepo.Save(job)

	var jobResult entity.Job
	err = suite.DB.QueryRow("SELECT id, output_bucket_path, status, video_id, created_at, updated_at FROM job WHERE id = ?", job.ID).
		Scan(&jobResult.ID, &jobResult.OutputBucketPath, &jobResult.Status, &jobResult.VideoID, &jobResult.CreatedAt, &jobResult.UpdatedAt)

	suite.NoError(err)
	suite.Equal(job.ID, jobResult.ID)
	suite.Equal(job.OutputBucketPath, jobResult.OutputBucketPath)
	suite.Equal(job.Status, jobResult.Status)
	suite.Equal(job.CreatedAt.UnixMilli(), jobResult.CreatedAt.UnixMilli())
	suite.Equal(job.UpdatedAt.UnixMilli(), jobResult.UpdatedAt.UnixMilli())
}

// func (suite *JobRepositoryTestSuite) TestGivenValidId_WhenFind_ThenShouldRetrieveTheJob() {
// 	repo := database.NewVideoRepository(suite.DB)

// 	id := uuid.NewString()

// 	video, _ := entity.NewVideo(id, "resource_id", "file_path")
// 	repo.Save(video)

// 	videoResult, err := repo.Find(id)

// 	suite.NoError(err)

// 	suite.Equal(video.ID, videoResult.ID)
// 	suite.Equal(video.ResourceID, videoResult.ResourceID)
// 	suite.Equal(video.FilePath, videoResult.FilePath)
// 	suite.Equal(video.CreatedAt.UnixMilli(), videoResult.CreatedAt.UnixMilli())
// }

// func (suite *VideoRepositoryTestSuite) TestGivenInValidId_WhenFind_ThenShouldNotRetrieveAnyJob() {
// 	repo := database.NewVideoRepository(suite.DB)

// 	id := uuid.NewString()

// 	video, _ := entity.NewVideo(id, "resource_id", "file_path")
// 	repo.Save(video)

// 	videoResult, err := repo.Find("123")

// 	suite.Error(err)
// 	suite.Nil(videoResult)
// }
