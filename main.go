package main

import (
	"fmt"
	"github.com/GeorgeMac/pontoon/build"
	"github.com/GeorgeMac/pontoon/config"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/fsouza/go-dockerclient"
	"os"
	"time"
)

var now func() time.Time = time.Now

func main() {
	opts := config.ParseOptions()

	client, err := docker.NewClient(opts.Host)
	if err != nil {
		panic(err)
	}

	fmt.Println("Fetching repository ", opts.Url)

	projects, err := project.NewGitProjects(opts.Dir)
	if err != nil {
		panic(err)
	}

	project, err := projects.Get(opts.Url)
	if err != nil {
		panic(err)
	}

	builder := build.NewBuilder(client, os.Stdout)
	if err := builder.BuildProject(project, "georgemac/hellos"); err != nil {
		panic(err)
	}

	image, err := client.InspectImage("georgemac/hellos")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", image)
}
