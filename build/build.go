package build

import (
	"bytes"
	"github.com/GeorgeMac/pontoon/project"
	"github.com/fsouza/go-dockerclient"
	"os"
)

type Builder struct {
	client *docker.Client
	out    *os.File
}

func NewBuilder(client *docker.Client, out *os.File) *Builder {
	return &Builder{client: client, out: out}
}

func (b *Builder) BuildProject(p project.Project, name string) error {
	input := bytes.NewBuffer(nil)
	if err := p.WriteTo(input); err != nil {
		return err
	}

	return b.client.BuildImage(docker.BuildImageOptions{
		Name:           name,
		InputStream:    input,
		OutputStream:   b.out,
		RmTmpContainer: true,
	})
}
