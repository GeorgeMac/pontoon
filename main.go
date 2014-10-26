package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/GeorgeMac/pontoon/build"
	"github.com/GeorgeMac/pontoon/client"
	"github.com/GeorgeMac/pontoon/config"
	"github.com/GeorgeMac/pontoon/jobs"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/GeorgeMac/pontoon/service"
)

var now func() time.Time = time.Now

func main() {
	var confpth string
	flag.StringVar(&confpth, "conf", "./config.yml", "Location of YAML config file")
	flag.Parse()

	conf, err := config.Parse(confpth)
	if err != nil {
		log.Fatal(err)
	}

	pconf, dconf := conf.Pontoon, conf.Docker

	client, err := client.New(dconf.Host, dconf.CaPem, dconf.CertPem, dconf.KeyPem)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(); err != nil {
		log.Fatal(err)
	}

	projects, err := project.NewGitProjects(pconf.Dir)
	if err != nil {
		panic(err)
	}

	// construct a build job factory
	factory := &build.BuildJobFactory{
		Client:   client,
		Projects: projects,
	}

	// construct and begin a queue for jobs
	queue := jobs.NewJobQueue(pconf.Workers)

	// create the new http service
	s := service.NewService(queue, factory)

	log.Println("Starting service port :8080")
	http.ListenAndServe(":8080", s)
}
