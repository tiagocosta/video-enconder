package event

import "time"

type VideoRequested struct {
	Name    string
	Payload interface{}
}

func NewVideoRequested() *VideoRequested {
	return &VideoRequested{
		Name: "VideoRequested",
	}
}

func (e *VideoRequested) GetName() string {
	return e.Name
}

func (e *VideoRequested) GetPayload() interface{} {
	return e.Payload
}

func (e *VideoRequested) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *VideoRequested) GetDateTime() time.Time {
	return time.Now()
}
