package config

import (
	"flag"
)

type Option struct {
	Host, Dir, Url, Author string
}

func ParseOptions() (o Option) {
	flag.StringVar(&o.Host, "h", "tcp://localhost:4243", "Docker host")
	flag.StringVar(&o.Dir, "d", "/var/pontoon", "Location where pontoon can clone projects")
	flag.StringVar(&o.Author, "author", "georgemac", "Author of the images to be built")
	flag.Parse()
	o.Url = flag.Arg(0)
	return
}
