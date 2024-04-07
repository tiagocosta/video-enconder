package database_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) SetupAllSuite() {
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

func (suite *RepositoryTestSuite) TearDownSuite() {
	suite.DB.Close()
}
