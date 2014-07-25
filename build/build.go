package build

import (
	"bytes"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/fsouza/go-dockerclient"
	"io"
)

type BuildJobFactory struct {
	Client   *docker.Client
	Projects project.Projects
}

func (b *BuildJobFactory) NewJob(name, url string) (bj *BuildJob, err error) {
	var p project.Project
	if p, err = b.Projects.Get(url); err != nil {
		return
	}

	return &BuildJob{
		name:    name,
		client:  b.Client,
		project: p,
	}, nil
}

type BuildJob struct {
	name    string
	project project.Project
	client  *docker.Client
	out     io.Writer
}

func (b *BuildJob) Run() error {
	input := bytes.NewBuffer(nil)
	if err := b.project.WriteTo(input); err != nil {
		return err
	}

	return b.client.BuildImage(docker.BuildImageOptions{
		Name:           b.name,
		InputStream:    input,
		OutputStream:   b.out,
		RmTmpContainer: true,
		NoCache:        true,
	})
}

func (b *BuildJob) SetOutput(wr io.Writer) {
	b.out = wr
}

func (b *BuildJob) String() string {
	return b.name
}
