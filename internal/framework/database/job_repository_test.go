package database_test

import (
	"github.com/google/uuid"
	"github.com/tiagocosta/video-enconder/internal/entity"
	"github.com/tiagocosta/video-enconder/internal/framework/database"

	_ "github.com/mattn/go-sqlite3"
)

type JobRepositoryTestSuite struct {
	RepositoryTestSuite
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
		Scan(&job.ID, &job.OutputBucketPath, &job.Status, &job.VideoID, &job.CreatedAt, &job.UpdatedAt)

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
