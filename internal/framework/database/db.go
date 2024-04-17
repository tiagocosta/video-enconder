package database

import (
	"database/sql"
	"fmt"

	"github.com/tiagocosta/video-enconder/configs"
)

func SqlDB() *sql.DB {
	configs := configs.Config()

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(172.26.0.1:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	return db
}
