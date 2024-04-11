package database

import (
	"database/sql"
	"time"

	"github.com/tiagocosta/video-enconder/internal/entity"
)

type JobRepository struct {
	DB *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{DB: db}
}

func (r *JobRepository) Save(job *entity.Job) error {
	stmt, err := r.DB.Prepare("INSERT INTO job (id, output_bucket_path, status, video_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(job.ID, job.OutputBucketPath, job.Status, job.VideoID, job.CreatedAt, job.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *JobRepository) Find(id string) (*entity.Job, error) {
	var job entity.Job
	err := r.DB.QueryRow("SELECT id, output_bucket_path, status, video_id, created_at, updated_at FROM job WHERE id = ?", id).
		Scan(&job.ID, &job.OutputBucketPath, &job.Status, &job.VideoID, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}
	videoRepo := NewVideoRepository(r.DB)
	video, err := videoRepo.Find(job.VideoID)
	if err != nil {
		return nil, err
	}
	job.Video = video

	return &job, nil
}

func (r *JobRepository) Update(job *entity.Job) error {
	stmt, err := r.DB.Prepare("UPDATE job SET output_bucket_path = ?, status = ?, video_id = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(job.OutputBucketPath, job.Status, job.VideoID, time.Now(), job.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *JobRepository) List(video *entity.Video) ([]entity.Job, error) {
	rows, err := r.DB.Query("SELECT id, output_bucket_path, status, video_id, created_at, updated_at FROM job WHERE video_id = ?", video.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := []entity.Job{}
	for rows.Next() {
		var id, output_bucket_path, status, video_id string
		var created_at, updated_at time.Time
		if err := rows.Scan(&id, &output_bucket_path, &status, &video_id, &created_at, &updated_at); err != nil {
			return nil, err
		}
		jobs = append(jobs, entity.Job{
			ID:               id,
			OutputBucketPath: output_bucket_path,
			Status:           entity.Status(status),
			VideoID:          video_id,
			Video:            video,
			CreatedAt:        created_at,
			UpdatedAt:        updated_at,
		})
	}
	return jobs, nil
}
