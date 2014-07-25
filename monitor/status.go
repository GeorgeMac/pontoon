package monitor

// Status is an enum, which describes the state of a Job
type Status int

const (
	UNKNOWN Status = iota
	PENDING
	ACTIVE
	COMPLETE
	FAILED
)

// Trackable defines the type of something which identifiable, has a
// status and a final output upon completion (Status() -> COMPLETE).
type Trackable interface {
	Status() Status
	Output() (string, error)
	Report() TrackableReport
}

type TrackableReport struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}
