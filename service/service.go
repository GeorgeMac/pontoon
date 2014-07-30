package service

import (
	"encoding/json"
	"fmt"
	"github.com/GeorgeMac/pontoon/build"
	"github.com/GeorgeMac/pontoon/jobs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const BUILD_TIMEOUT time.Duration = 2 * time.Minute

type Service struct {
	queue  *jobs.JobQueue
	fact   *build.BuildJobFactory
	store  *jobs.Store
	router *mux.Router
}

func NewService(queue *jobs.JobQueue, fact *build.BuildJobFactory) (s *Service) {
	s = &Service{
		queue:  queue,
		fact:   fact,
		router: mux.NewRouter(),
		store:  jobs.NewStore(),
	}
	s.router.Methods("POST").Subrouter().HandleFunc("/jobs", s.submit)
	s.router.Methods("GET").Subrouter().HandleFunc("/jobs", s.list)

	s.router.Methods("POST").Subrouter().HandleFunc("/jobs/{id}", s.build)
	s.router.Methods("GET").Subrouter().HandleFunc("/jobs/{id}", s.job)
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

func (s *Service) job(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(s.store.FullReport(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) list(w http.ResponseWriter, req *http.Request) {
	if err := json.NewEncoder(w).Encode(s.store.List()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) build(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	job, err := s.store.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&BuildResponse{
			Msg: fmt.Sprintf("Build with name %s is missing", id),
		})
	}

	// push the job in to the build queue
	s.queue.Push(job)
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

	if st := s.store.Report(request.Name); st.Status != "UNKNOWN" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&BuildResponse{
			Msg:    fmt.Sprintf("Build with name %s already exists", request.Name),
			Status: st.Status,
		})
	}

	if err := s.store.Put(request.Name, job); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}
