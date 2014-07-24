package jobs

import (
	"bytes"
	"fmt"
	"github.com/GeorgeMac/pontoon/monitor"
	"io"
	"sync"
	"time"
)

type Executor interface {
	Run() error
	SetOutput(io.Writer)
	fmt.Stringer
}

var Now func() time.Time = time.Now

// Job controls access to the state of a runner and its eventual
// output.
type Job struct {
	mu        *sync.RWMutex
	st        monitor.Status
	exec      Executor
	buffer    *bytes.Buffer
	output    string
	Signal    chan monitor.Status
	CreatedAt time.Time
}

func NewJob(exec Executor) *Job {
	when := Now()

	desc := "Build job for executor %s created at %v\n"
	desc = fmt.Sprintf(desc, exec, when)

	buffer := bytes.NewBufferString(desc)
	exec.SetOutput(buffer)

	return &Job{
		mu:        &sync.RWMutex{},
		exec:      exec,
		buffer:    buffer,
		Signal:    make(chan monitor.Status),
		CreatedAt: when,
	}
}

// Status obtains a read lock and then returns the Job status
func (j *Job) Status() (s monitor.Status) {
	j.mu.RLock()
	s = j.st
	j.mu.RUnlock()
	return
}

// Output returns the output string of the job, but only once the
// jobs status has a value of COMPLETE
func (j *Job) Output() (string, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	if j.st < monitor.COMPLETE {
		return "", JobIncompleteError{}
	}

	return j.output, nil
}

// begin sets the job status to active and returns the build.BuildJob
// and starts a routine to consume job output on completion.
func (j *Job) begin() (Executor, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.st > monitor.PENDING {
		return nil, JobBegunError{}
	}

	j.st = monitor.ACTIVE

	go j.finish()

	return j.exec, nil
}

// finish blocks until the job has signalled it is complete
// It then obtains a lock, sets its status and reads out the
// build jobs output
func (j *Job) finish() {
	st := <-j.Signal
	j.mu.Lock()
	j.output = j.buffer.String()
	j.st = st
	j.mu.Unlock()
}

// JobIncompleteError is called when output is requested
// before job completion
type JobIncompleteError struct{}

func (j JobIncompleteError) Error() string {
	return "Job currently active"
}

// JobBegunError is fired when begin is called on a job
// that has begun
type JobBegunError struct{}

func (j JobBegunError) Error() string {
	return "Job has already begun"
}
