package build

import (
	"bytes"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/fsouza/go-dockerclient"
	"io"
)

type Builder struct {
	client   *docker.Client
	projects project.Projects
}

func NewBuilder(client *docker.Client, projects project.Projects) *Builder {
	return &Builder{
		client:   client,
		projects: projects,
	}
}

func (b *Builder) BuildProject(job *BuildJob) error {
	defer close(job.Done)

	p, err := b.projects.Get(job.Url)
	if err != nil {
		return err
	}

	input := bytes.NewBuffer(nil)
	if err := p.WriteTo(input); err != nil {
		return err
	}

	return b.client.BuildImage(docker.BuildImageOptions{
		Name:           n,
		InputStream:    input,
		OutputStream:   o,
		RmTmpContainer: true,
	})
}

type BuildJob struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Out  io.Writer
	Done chan struct{}
}

func NewBuildJob(url, name string) *BuildJob {
	return &BuildJob{
		Url:  url,
		Name: name,
		Done: make(chan struct{}),
	}
}
