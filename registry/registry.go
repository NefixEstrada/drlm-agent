package registry

import (
	"context"
	"encoding/json"

	"github.com/brainupdaters/drlm-common/pkg/os"
	"go.etcd.io/etcd/clientv3"
)

const base = "/registry/drlm-agent/"

type agent struct {
	Addr     string        `json:"addr,omitempty"`
	Metadata agentMetadata `json:"metadata,omitempty"`
}

type agentMetadata struct {
	Version       string  `json:"version,omitempty"`
	Arch          os.Arch `json:"arch,omitempty"`
	OS            os.OS   `json:"os,omitempty"`
	OSVersion     string  `json:"os_version,omitempty"`
	Distro        string  `json:"distro,omitempty"`
	DistroVersion string  `json:"distro_version,omitempty"`
}

func Register(urls []string, uuid string, addr string) {
	// client.New
	cli, err := clientv3.NewFromURLs(urls)
	if err != nil {
		panic(err)
	}

	a := &agent{
		Addr:     addr,
		Metadata: agentMetadata{},
	}

	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	cli.Put(context.TODO(), base+uuid, string(b))

	// cli, cerr := cv.NewFromURL("http://localhost:2379")
	// r := &etcdnaming.GRPCResolver{Client: cli}
	// b := grpc.RoundRobin(r)
	// conn, gerr := grpc.Dial("my-service", grpc.WithBalancer(b), grpc.WithBlock())
}
