package database

import (
	"database/sql"

	"github.com/tiagocosta/video-enconder/internal/entity"
)

type VideoRepository struct {
	DB *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{DB: db}
}

func (r *VideoRepository) Save(video *entity.Video) error {
	stmt, err := r.DB.Prepare("INSERT INTO video (id, resource_id, file_path, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(video.ID, video.ResourceID, video.FilePath, video.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *VideoRepository) Find(id string) (*entity.Video, error) {
	var video entity.Video
	err := r.DB.QueryRow("SELECT id, resource_id, file_path, created_at FROM video WHERE id = ?", id).
		Scan(&video.ID, &video.ResourceID, &video.FilePath, &video.CreatedAt)
	if err != nil {
		return nil, err
	}
	jobRepo := NewJobRepository(r.DB)
	jobs, err := jobRepo.List(&video)
	if err != nil {
		return nil, err
	}
	video.Jobs = jobs

	return &video, nil
}
