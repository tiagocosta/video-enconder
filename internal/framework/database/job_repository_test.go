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
	suite.Equal(job.VideoID, jobResult.VideoID)
	suite.Equal(job.CreatedAt.UnixMilli(), jobResult.CreatedAt.UnixMilli())
	suite.Equal(job.UpdatedAt.UnixMilli(), jobResult.UpdatedAt.UnixMilli())
}

func (suite *JobRepositoryTestSuite) TestGivenValidId_WhenFind_ThenShouldRetrieveTheJobAndAssociatedVideo() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "file_path")
	videoRepo := database.NewVideoRepository(suite.DB)
	videoRepo.Save(video)

	job, _ := entity.NewJob("output_path", "Pending", video)
	jobRepo := database.NewJobRepository(suite.DB)
	jobRepo.Save(job)

	jobResult, err := jobRepo.Find(job.ID)

	suite.NoError(err)
	suite.Equal(job.ID, jobResult.ID)
	suite.Equal(job.OutputBucketPath, jobResult.OutputBucketPath)
	suite.Equal(job.Status, jobResult.Status)
	suite.Equal(job.VideoID, jobResult.VideoID)
	suite.Equal(job.VideoID, jobResult.Video.ID)
	suite.Equal(job.CreatedAt.UnixMilli(), jobResult.CreatedAt.UnixMilli())
	suite.Equal(job.UpdatedAt.UnixMilli(), jobResult.UpdatedAt.UnixMilli())
}

func (suite *JobRepositoryTestSuite) TestGivenInvalidId_WhenFind_ThenShouldNotRetrieveAnyJob() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "file_path")
	videoRepo := database.NewVideoRepository(suite.DB)
	videoRepo.Save(video)

	job, _ := entity.NewJob("output_path", "Pending", video)
	jobRepo := database.NewJobRepository(suite.DB)
	jobRepo.Save(job)

	jobResult, err := jobRepo.Find("123")
	suite.Error(err)
	suite.Nil(jobResult)
}

func (suite *JobRepositoryTestSuite) TestGivenValidVideoId_WhenList_ThenShouldRetrieveAllJobsAssociatedToThatVideo() {
	video, _ := entity.NewVideo(uuid.NewString(), "resource_id", "file_path")
	videoRepo := database.NewVideoRepository(suite.DB)
	videoRepo.Save(video)

	job1, _ := entity.NewJob("output_path", "Pending", video)
	job2, _ := entity.NewJob("output_path", "Pending", video)
	job3, _ := entity.NewJob("output_path", "Pending", video)
	jobRepo := database.NewJobRepository(suite.DB)
	jobRepo.Save(job1)
	jobRepo.Save(job2)
	jobRepo.Save(job3)

	jobsResult, err := jobRepo.List(video)
	suite.NoError(err)
	suite.Len(jobsResult, 3)

	suite.Equal(job1.ID, jobsResult[0].ID)
	suite.Equal(job1.OutputBucketPath, jobsResult[0].OutputBucketPath)
	suite.Equal(job1.Status, jobsResult[0].Status)
	suite.Equal(job1.VideoID, jobsResult[0].VideoID)
	suite.Equal(job1.VideoID, jobsResult[0].Video.ID)
	suite.Equal(job1.CreatedAt.UnixMilli(), jobsResult[0].CreatedAt.UnixMilli())
	suite.Equal(job1.UpdatedAt.UnixMilli(), jobsResult[0].UpdatedAt.UnixMilli())

	suite.Equal(job2.ID, jobsResult[1].ID)
	suite.Equal(job2.OutputBucketPath, jobsResult[1].OutputBucketPath)
	suite.Equal(job2.Status, jobsResult[1].Status)
	suite.Equal(job2.VideoID, jobsResult[1].VideoID)
	suite.Equal(job2.VideoID, jobsResult[1].Video.ID)
	suite.Equal(job2.CreatedAt.UnixMilli(), jobsResult[1].CreatedAt.UnixMilli())
	suite.Equal(job2.UpdatedAt.UnixMilli(), jobsResult[1].UpdatedAt.UnixMilli())

	suite.Equal(job3.ID, jobsResult[2].ID)
	suite.Equal(job3.OutputBucketPath, jobsResult[2].OutputBucketPath)
	suite.Equal(job3.Status, jobsResult[2].Status)
	suite.Equal(job3.VideoID, jobsResult[2].VideoID)
	suite.Equal(job3.VideoID, jobsResult[2].Video.ID)
	suite.Equal(job3.CreatedAt.UnixMilli(), jobsResult[2].CreatedAt.UnixMilli())
	suite.Equal(job3.UpdatedAt.UnixMilli(), jobsResult[2].UpdatedAt.UnixMilli())
}
