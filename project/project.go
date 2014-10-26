package project

import (
	"io"

	"github.com/fsouza/go-dockerclient"
)

type BuildOptions docker.BuildImageOptions

type Project interface {
	WriteTo(io.Writer) error
	Pull() error
	Ref() (string, error)
}
