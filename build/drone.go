package build

const defaultPipeline = `
kind: pipeline
type: exec
name: default

platform:
  os: linux
  arch: amd64

workspace: 
  path: /root/go/src

steps:
- name: build
  commands:
  - export GOROOT=/usr/local/go
  - export GOPATH=/root/go
  - export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
  - go mod init
  - go mod vendor
  - go build -o drone-test-go -v
  - mkdir -p /data/rolling-build/drone-test-go
  - cp -a drone-test-go /data/rolling-build/drone-test-go/
  - cd /data/rolling-build/drone-test-go/
  - echo "docker build . ${PWD}"
`

type Drone struct {
	Addr 	string
}

func NewDroneServer(Addr string) *Drone {
	return &Drone{
		Addr: addr,
	}
}