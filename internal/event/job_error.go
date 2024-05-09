package event

import "time"

type JobError struct {
	Name    string
	Payload interface{}
}

func NewJobError() *JobError {
	return &JobError{
		Name: "JobError",
	}
}

func (e *JobError) GetName() string {
	return e.Name
}

func (e *JobError) GetPayload() interface{} {
	return e.Payload
}

func (e *JobError) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *JobError) GetDateTime() time.Time {
	return time.Now()
}
