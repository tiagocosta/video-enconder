package entity

type VideoRepositoryInterface interface {
	Save(video *Video) error
	Find(id string) (*Video, error)
}
