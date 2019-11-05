export GO111MODULE=on
BINARY_NAME=drlm-agent

all: test build

test:
	go test -cover ./...

build:
	go build -o $(BINARY_NAME) drlm-agent.go