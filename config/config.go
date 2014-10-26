package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v1"
)

type Config struct {
	Pontoon PontoonConf `yaml:"pontoon"`
	Docker  DockerConf  `yaml:"docker"`
}

type PontoonConf struct {
	Dir     string `yaml:"directory"`
	Workers int    `yaml:"workers"`
}

type DockerConf struct {
	Host    string `yaml:"host"`
	CaPem   string `yaml:"ca_pem"`
	CertPem string `yaml:"cert_pem"`
	KeyPem  string `yaml:"key_pem"`
}

func Parse(pth string) (c Config, err error) {
	data, err := ioutil.ReadFile(pth)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(data, &c); err != nil {
		return
	}

	if c.Pontoon.Workers == 0 {
		c.Pontoon.Workers = 1
	}
	return
}
