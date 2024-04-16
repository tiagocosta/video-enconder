package event

import "time"

type JobCompleted struct {
	Name    string
	Payload interface{}
}

func NewJobCompleted() *JobCompleted {
	return &JobCompleted{
		Name: "JobCompleted",
	}
}

func (e *JobCompleted) GetName() string {
	return e.Name
}

func (e *JobCompleted) GetPayload() interface{} {
	return e.Payload
}

func (e *JobCompleted) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *JobCompleted) GetDateTime() time.Time {
	return time.Now()
}
