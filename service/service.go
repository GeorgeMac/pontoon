package service

import (
	"encoding/json"
	"fmt"
	"github.com/GeorgeMac/pontoon/build"
	"github.com/GeorgeMac/pontoon/jobs"
	"github.com/GeorgeMac/pontoon/monitor"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const BUILD_TIMEOUT time.Duration = 2 * time.Minute

type Service struct {
	queue  *jobs.JobQueue
	fact   *build.BuildJobFactory
	mont   *monitor.Monitor
	router *mux.Router
}

func NewService(queue *jobs.JobQueue, fact *build.BuildJobFactory) (s *Service) {
	s = &Service{
		queue:  queue,
		fact:   fact,
		router: mux.NewRouter(),
		mont:   monitor.NewMonitor(),
	}
	s.router.Methods("POST").Subrouter().HandleFunc("/jobs", s.submit)
	s.router.Methods("Get").Subrouter().HandleFunc("/jobs", s.list)
	return
}

func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

type BuildRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type BuildResponse struct {
	Msg    string `json:"message"`
	Status string `json:"status"`
}

func (s *Service) list(w http.ResponseWriter, req *http.Request) {
	if err := json.NewEncoder(w).Encode(s.mont.List()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func (s *Service) submit(w http.ResponseWriter, req *http.Request) {
	request := BuildRequest{}
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	// get a build job from the factory
	bj, err := s.fact.NewJob(request.Name, request.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// wrap in a jobs.Job
	job := jobs.NewJob(bj)

	if st := s.mont.Report(request.Name); st.Status > monitor.UNKNOWN {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&BuildResponse{
			Msg:    fmt.Sprintf("Build with name %s already exists", request.Name),
			Status: fmt.Sprintf("%d", st.Status),
		})
	}

	if err := s.mont.Put(request.Name, job); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// push the job in to the build queue
	s.queue.Push(job)
}
