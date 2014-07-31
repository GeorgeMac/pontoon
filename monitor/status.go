package monitor

// Status is an enum, which describes the state of a Job
type Status int

const (
	UNKNOWN Status = iota
	READY
	PENDING
	ACTIVE
	COMPLETE
	FAILED
	ABORTED
)

func (s Status) String() string {
	switch s {
	case UNKNOWN:
		return "UNKNOWN"
	case READY:
		return "READY"
	case PENDING:
		return "PENDING"
	case ACTIVE:
		return "ACTIVE"
	case COMPLETE:
		return "COMPLETE"
	case FAILED:
		return "FAILED"
	case ABORTED:
		return "ABORTED"
	default:
		return "UNEXPECTED STATUS"
	}
}

// Reportable can produce a Report struct. This report should have a
// status < COMPLETE. It can also produce a slice of previous Completed.
type Reportable interface {
	Report() Report
	History() History
}

type FullReport struct {
	Report
	History History `json:"history"`
}

type Report struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type History []CompletedReport

type CompletedReport struct {
	Status string `json:"status"`
	Id     int16  `json:"id"`
	Output string `json:"output"`
}
