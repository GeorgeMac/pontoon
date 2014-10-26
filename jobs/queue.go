package jobs

import (
	"errors"
	"log"

	"github.com/GeorgeMac/pontoon/monitor"
)

var NotEnoughWorkersError = errors.New("Not Enough Workers")

type JobQueue struct {
	workers int
	jobs    chan *Job
	stop    chan struct{}
}

func NewJobQueue(workers int) (*JobQueue, error) {
	if workers <= 0 {
		return nil, NotEnoughWorkersError
	}

	q := &JobQueue{
		workers: workers,
		jobs:    make(chan *Job, workers),
		stop:    make(chan struct{}),
	}

	for i := 0; i < workers; i++ {
		go q.begin()
	}

	return q, nil
}

func (q *JobQueue) begin() {
	for {
		select {
		case job := <-q.jobs:
			exec, err := job.begin()
			if err != nil {
				log.Println(err.Error())
				if _, ok := err.(JobBegunError); ok {
					continue
				}
				job.Signal <- monitor.FAILED
				continue
			}

			if err = exec.Run(); err != nil {
				log.Println(err.Error())
				job.Signal <- monitor.FAILED
				continue
			}

			job.Signal <- monitor.COMPLETE
		case <-q.stop:
			return
		}
	}
}

func (q *JobQueue) Push(job *Job) {
	job.setStatus(monitor.PENDING)
	q.jobs <- job
}

func (q *JobQueue) Stop() {
	for i := 0; i < q.workers; i++ {
		q.stop <- struct{}{}
	}
	close(q.stop)
}
