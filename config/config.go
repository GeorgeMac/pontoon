package config

import (
	"flag"
)

type Option struct {
	Host, Dir string
}

func ParseOptions() (o Option) {
	flag.StringVar(&o.Host, "h", "tcp://localhost:4243", "Docker host")
	flag.StringVar(&o.Dir, "d", "/var/pontoon", "Location where pontoon can clone projects")
	flag.Parse()
	return
}
