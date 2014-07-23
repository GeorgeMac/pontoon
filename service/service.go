package service

import (
	"encoding/json"
	"github.com/GeorgeMac/pontoon/build"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const BUILD_TIMEOUT time.Duration = 2 * time.Minute

type Service struct {
	queue  *build.BuildQueue
	router *mux.Router
}

func NewService(queue *build.BuildQueue) (s *Service) {
	s = &Service{
		queue:  queue,
		router: mux.NewRouter(),
	}
	s.router.Methods("POST").Subrouter().HandleFunc("/jobs", s.submit)
	return
}

func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Service) submit(w http.ResponseWriter, req *http.Request) {
	job := build.NewBuildJob("", "", w)
	if err := json.NewDecoder(req.Body).Decode(&job); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// push the job in to the build queue
	s.queue.Push(job)

	for {
		select {
		case <-job.Done:
			return
		case <-time.After(BUILD_TIMEOUT):
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
	}
}
