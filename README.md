# Pontoon

#### A Docker build and test server

Very early stage development git project -> docker image service.
Currently only has one test. This is mostly just an experiment!

##### Current Capabilities
``` pontoon -d “/var/pontoon” -h “tcp://localhost:4321” ```

1. Pulls git projects down in to `/var/pontoon`
2. Connects to docker at host `tcp://localhost:4321`
3. Post a build job to `http://localhost:8080/jobs` in the following JSON format:

```json
{
	“name”:”georgemac/demo”,
	“url”:”http://github.com/GeorgeMac/hello-server”
}
```
