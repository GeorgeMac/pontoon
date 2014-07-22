package project

import (
	"github.com/fsouza/go-dockerclient"
	"io"
)

type BuildOptions docker.BuildImageOptions

type Projects interface {
	Get(string) (Project, error)
}

type Project interface {
	WriteTo(io.Writer) error
	Pull() error
	Ref() (string, error)
}
