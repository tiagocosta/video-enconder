package database

import (
	"log"

	"github.com/tiagocosta/video-enconder/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB            *gorm.DB
	Dsn           string
	DsnTest       string
	Debug         bool
	AutoMigrateDB bool
	Env           string
}

func NewDB() *Database {
	return &Database{}
}

func NewDBTest() *gorm.DB {
	dbInstance := NewDB()
	dbInstance.Env = "Test"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDB = true
	dbInstance.Debug = true

	conn, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("test db error: %v", err)
	}

	return conn
}

func (db *Database) Connect() (*gorm.DB, error) {
	var err error

	if db.Env == "Test" {
		db.DB, err = gorm.Open(sqlite.Open(db.DsnTest), &gorm.Config{})
	} else {
		db.DB, err = gorm.Open(postgres.Open(db.Dsn), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	if db.Debug {
		db.DB.Logger.LogMode(logger.Error)
	}

	if db.AutoMigrateDB {
		db.DB.AutoMigrate(&entity.Video{}, &entity.Job{})
	}

	return db.DB, nil
}
