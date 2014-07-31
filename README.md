# Pontoon

#### A Docker build and test server

Very early stage development git project -> docker image service.
Currently only has one test. This is mostly just an experiment!

##### Current Capabilities
``` pontoon -d "/var/pontoon" -h "tcp://localhost:4321" ```

- Pulls git projects down in to `/var/pontoon`.
- Connects to docker at host `tcp://localhost:4321`.

##### API

```
GET  /jobs - List of current build jobs that can be built
```
Get a list of jobs and their status currently retained in memory.

```bash
curl localhost:8080/jobs
```
for example might return:
```json
[{"name":"georgemac", "status":"READY"}]
```
---
```
POST /jobs - Send a new job
```
Name and job and attribute it to a git url. Readying it for building.

```bash
curl -X POST -d http://localhost:8080/jobs 
```
for example would return:
```json
{
	"name":"georgemac",
	"url":"http://github.com/GeorgeMac/hello-server"
}
```
---
```
GET  /jobs/{id} - Full report of job with {id} and its previous builds
```
Get a summary of a job, including it's previous builds.

###### example
``` bash
curl http://localhost:8080/job/georgemac
```
###### response
```json
{
  "status": "READY",
  "name": "georgemac",
  "Previous": [
    {
      "output": "Build job for executor georgemac created at 2014-07-30 18:55:13.46836866 +0100 BST\nStep 0 : FROM google/golang\n ---> fa77fdfe2188\nStep 1 : MAINTAINER George MacRorie github.com/GeorgeMac\n ---> Running in 42c931c2922b\n ---> fdd04f4eb080\nStep 2 : WORKDIR /gopath/src/hello-server\n ---> Running in e7e3305fffec\n ---> f0f542ed4a20\nStep 3 : ADD . /gopath/src/hello-server/\n ---> 4d9d0620a581\nStep 4 : RUN pwd\n ---> Running in f43cd7eee77b\n/gopath/src/hello-server\n ---> a5597e5995df\nStep 5 : RUN go test ./...\n ---> Running in 8ea120200e0d\nok  \thello-server\t0.003s\n ---> 1bb64e351798\nStep 6 : RUN go install ./...\n ---> Running in 0b0290b9e1e8\n ---> dba44cc36663\nStep 7 : CMD []\n ---> Running in e1459201db2c\n ---> 8a7750276d0d\nStep 8 : ENTRYPOINT [\"/gopath/bin/hello-server\"]\n ---> Running in 27773dea008f\n ---> 4d973cdb40eb\nSuccessfully built 4d973cdb40eb\nRemoving intermediate container 42c931c2922b\nRemoving intermediate container e7e3305fffec\nRemoving intermediate container a339c7e1858f\nRemoving intermediate container f43cd7eee77b\nRemoving intermediate container 8ea120200e0d\nRemoving intermediate container 0b0290b9e1e8\nRemoving intermediate container e1459201db2c\nRemoving intermediate container 27773dea008f\n",
      "id": 1,
      "status": "COMPLETE"
    }
  ]
}
```

---
```
POST /jobs/{id} - Trigger a new build for job with {id}
```
###### example
```bash
curl -X POST http://localhost:8080/job/georgemac
```
This will trigger a build of the job with id `georgemac`

Again remember this is all experimental and very much untested. Keep watching!
