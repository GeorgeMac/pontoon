package build

import (
	"bytes"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/fsouza/go-dockerclient"
	"io"
	"log"
)

type BuildQueue struct {
	builder  *Builder
	projects project.Projects
	jobs     chan *BuildJob
	stop     chan struct{}
}

func NewBuilderQueue(bld *Builder, prj project.Projects) *BuildQueue {
	return &BuildQueue{
		builder:  bld,
		projects: prj,
		jobs:     make(chan *BuildJob),
		stop:     make(chan struct{}),
	}
}

func (b *BuildQueue) Begin() {
	for {
		select {
		case job := <-b.jobs:
			defer close(job.Done)
			p, err := b.projects.Get(job.Url)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if err = b.builder.BuildProject(p, job.Name, job.Out); err != nil {
				log.Println(err.Error())
				continue
			}
		case <-b.stop:
			return
		}
	}
}

func (b *BuildQueue) Push(job *BuildJob) {
	go func() {
		b.jobs <- job
	}()
}

func (b *BuildQueue) Stop() {
	b.stop <- struct{}{}
	close(b.stop)
}

type BuildJob struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Out  io.Writer
	Done chan struct{}
}

func NewBuildJob(url, name string, out io.Writer) *BuildJob {
	return &BuildJob{
		Url:  url,
		Name: name,
		Out:  out,
		Done: make(chan struct{}),
	}
}

type Builder struct {
	client *docker.Client
}

func NewBuilder(client *docker.Client) *Builder {
	return &Builder{client: client}
}

func (b *Builder) BuildProject(p project.Project, n string, o io.Writer) error {
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
