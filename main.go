package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/GeorgeMac/pontoon/archive"
	"github.com/fsouza/go-dockerclient"
	"os"
	"time"
)

var now func() time.Time = time.Now

func main() {
	client, err := docker.NewClient("tcp://localhost:4243")
	if err != nil {
		panic(err)
	}

	input := bytes.NewBuffer(nil)
	tr := tar.NewWriter(input)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := archive.Dir(dir, tr); err != nil {
		panic(err)
	}

	tr.Close()

	if err := client.BuildImage(docker.BuildImageOptions{
		Name:           "georgemac/demo",
		InputStream:    input,
		OutputStream:   os.Stdout,
		RmTmpContainer: true,
	}); err != nil {
		panic(err)
	}

	images, err := client.ListImages(false)
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Printf("%+v\n", image)
	}
}
