package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/fsouza/go-dockerclient"
)

var desired docker.ApiVersion

func init() {
	var err error
	if desired, err = docker.NewApiVersion("1.3.0"); err != nil {
		panic(err)
	}
}

func New(endp, ca, cert, key string) (*docker.Client, error) {
	roots := x509.NewCertPool()
	pemData, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	//add to pool
	roots.AppendCertsFromPEM(pemData)

	//create certificate
	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	//creates the new tls configuration using both the authority and certificate
	conf := &tls.Config{
		RootCAs:      roots,
		Certificates: []tls.Certificate{crt},
	}

	//create our own transport
	tr := &http.Transport{
		TLSClientConfig: conf,
	}

	host, err := url.Parse(endp)
	if err != nil {
		return nil, err
	}

	host.Scheme = "https"

	c, err := docker.NewClient(host.String())
	if err != nil {
		return nil, err
	}

	//create a new http client and set on our dockerclient
	c.HTTPClient = &http.Client{Transport: tr}

	env, err := c.Version()
	if err != nil {
		return nil, err
	}

	apiv := env.Get("Version")
	if apiv == "" {
		return nil, fmt.Errorf("Cannot Find Version in Environment")
	}

	obtained, err := docker.NewApiVersion(apiv)
	if err != nil {
		return nil, err
	}

	if desired.GreaterThan(obtained) {
		return nil, fmt.Errorf("API version %s not support", apiv)
	}

	return c, nil
}
