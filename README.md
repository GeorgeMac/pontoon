# Pontoon

#### A Docker build and test server

Very early stage development git project -> docker image service.
Currently only has one test. This is mostly just an experiment!

##### Current Capabilities
``` pontoon -d “/var/pontoon” -h “tcp://localhost:4321” ```

1. Pulls git projects down in to `/var/pontoon`.
2. Connects to docker at host `tcp://localhost:4321`.
3. Post a build job to `http://localhost:8080/jobs` in the following JSON format:

	```json
	{
		“name”:”georgemac/demo”,
		“url”:”http://github.com/GeorgeMac/hello-server”
	}
	```
4. Get a list of jobs and their status currently retained in memory.

	```bash
	curl localhost:8080/jobs
	```
	for example might return:
	```json
	[{"name":"georgemac/demo", "status":2}]
	```

	Until I stringify those status see `monitor.Status` enum. Currently they correspond as follows:
	```
	0: UNKNOWN
	1: PENDING
	2: ACTIVE
	3: COMPLETE
	4: FAILED
	```
	These are likely to change in the coming future.

Again remember this is all experimental and very much untested. By keep watching!
