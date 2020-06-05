export GO111MODULE=on
BINARY_NAME=drlm-agent

all: test build docker

test:
	go test -cover ./...

build:
	go build -o $(BINARY_NAME) drlm-agent.go

docker:
	docker build -t drlm-agent:1.0.0 .