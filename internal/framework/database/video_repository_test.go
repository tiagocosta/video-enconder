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

type VideoRepositoryTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func (suite *VideoRepositoryTestSuite) SetupSuite() {
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

func (suite *VideoRepositoryTestSuite) TearDownSuite() {
	suite.DB.Close()
}

func TestVideoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(VideoRepositoryTestSuite))
}

// func TestSuite(t *testing.T) {
// 	suite.Run(t, new(VideoRepositoryTestSuite))
// }

func (suite *VideoRepositoryTestSuite) TestGivenAnVideo_WhenSave_ThenShouldSaveVideo() {
	video, err := entity.NewVideo(uuid.NewString(), "resource_id", "file_path")
	suite.NoError(err)
	suite.NotNil(video)
	repo := database.NewVideoRepository(suite.DB)
	err = repo.Save(video)
	suite.NoError(err)

	var videoResult entity.Video
	err = suite.DB.QueryRow("SELECT id, resource_id, file_path, created_at FROM video WHERE id = ?", video.ID).
		Scan(&videoResult.ID, &videoResult.ResourceID, &videoResult.FilePath, &videoResult.CreatedAt)

	suite.NoError(err)
	suite.Equal(video.ID, videoResult.ID)
	suite.Equal(video.ResourceID, videoResult.ResourceID)
	suite.Equal(video.FilePath, videoResult.FilePath)
	suite.Equal(video.CreatedAt.UnixMilli(), videoResult.CreatedAt.UnixMilli())
}

func (suite *VideoRepositoryTestSuite) TestGivenValidId_WhenFind_ThenShouldRetrieveTheVideo() {
	repo := database.NewVideoRepository(suite.DB)

	id := uuid.NewString()

	video, _ := entity.NewVideo(id, "resource_id", "file_path")
	repo.Save(video)

	videoResult, err := repo.Find(id)

	suite.NoError(err)

	suite.Equal(video.ID, videoResult.ID)
	suite.Equal(video.ResourceID, videoResult.ResourceID)
	suite.Equal(video.FilePath, videoResult.FilePath)
	suite.Equal(video.CreatedAt.UnixMilli(), videoResult.CreatedAt.UnixMilli())
}

func (suite *VideoRepositoryTestSuite) TestGivenInValidId_WhenFind_ThenShouldNotRetrieveAnyVideo() {
	repo := database.NewVideoRepository(suite.DB)

	id := uuid.NewString()

	video, _ := entity.NewVideo(id, "resource_id", "file_path")
	repo.Save(video)

	videoResult, err := repo.Find("123")

	suite.Error(err)
	suite.Nil(videoResult)
}
