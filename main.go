package main

import (
	"github.com/GeorgeMac/pontoon/build"
	"github.com/GeorgeMac/pontoon/config"
	"github.com/GeorgeMac/pontoon/jobs"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/GeorgeMac/pontoon/service"
	"github.com/fsouza/go-dockerclient"
	"log"
	"net/http"
	"time"
)

var now func() time.Time = time.Now

func main() {
	opts := config.ParseOptions()

	client, err := docker.NewClient(opts.Host)
	if err != nil {
		panic(err)
	}

	projects, err := project.NewGitProjects(opts.Dir)
	if err != nil {
		panic(err)
	}

	// construct a build job factory
	factory := &build.BuildJobFactory{
		Client:   client,
		Projects: projects,
	}

	// construct and begin a queue for jobs
	queue := jobs.NewJobQueue(1)

	// create the new http service
	s := service.NewService(queue, factory)

	log.Println("Starting service port :8080")
	http.ListenAndServe(":8080", s)
}
