package main

import (
	"fmt"
	"github.com/GeorgeMac/pontoon/build"
	"github.com/GeorgeMac/pontoon/config"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/GeorgeMac/pontoon/service"
	"github.com/fsouza/go-dockerclient"
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

	// construct a project builder for the client
	builder := build.NewBuilder(client)

	// construct and begin a queue for the builder
	queue := build.NewBuilderQueue(builder, projects)
	go queue.Begin()

	s := service.NewService(queue)

	fmt.Println("Starting service port :8080")
	http.ListenAndServe(":8080", s)
}
