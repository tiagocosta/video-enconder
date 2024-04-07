package entity

type VideoRepositoryInterface interface {
	Save(video *Video) error
	Find(id string) (*Video, error)
}

type JobRepositoryInterface interface {
	Save(job *Job) error
	Find(id string) (*Job, error)
	Update(job *Job) error
	List(video *Video) []Job
}
