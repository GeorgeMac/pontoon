package monitor

// Status is an enum, which describes the state of a Job
type Status int

const (
	PENDING Status = iota
	ACTIVE
	COMPLETE
	FAILED
)

// Trackable defines the type of something which identifiable, has a
// status and a final output upon completion (Status() -> COMPLETE).
type Trackable interface {
	Status() Status
	Output() (string, error)
}
