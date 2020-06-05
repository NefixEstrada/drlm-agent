// SPDX-License-Identifier: AGPL-3.0-only

package grpc

import (
	"context"

	"github.com/brainupdaters/drlm-agent/binary"
	"github.com/brainupdaters/drlm-agent/cfg"
	"github.com/brainupdaters/drlm-agent/job"
	"github.com/brainupdaters/drlm-agent/models"

	"github.com/brainupdaters/drlm-common/pkg/core"
	drlm "github.com/brainupdaters/drlm-common/pkg/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"google.golang.org/grpc/metadata"
)

// API is the API version of the client
const API string = "v1.0.0"

var agentConn drlm.DRLM_AgentConnectionClient

// Init initializes the DRLM Core client
func Init(ctx context.Context, fs afero.Fs, join bool) {
	client, conn := core.NewClient(fs, cfg.Config.Core.TLS, cfg.Config.Core.CertPath, cfg.Config.Core.Host, cfg.Config.Core.Port)

	var err error
	agentConn, err = client.AgentConnection(prepareCtx(ctx))
	if err != nil {
		// TODO: RETRY
		log.Fatalf("error connecting to the core: %v", err)
	}

	if join {
		if err := agentConn.Send(&drlm.AgentConnectionFromAgent{
			MessageType: drlm.AgentConnectionFromAgent_MESSAGE_TYPE_JOIN_REQUEST,
			JoinRequest: &drlm.AgentConnectionFromAgent_JoinRequest{
				// TODO: Get the actual things
				Arch: drlm.Arch_ARCH_UNKNOWN,
				Os:   drlm.OS_OS_UNKNOWN,
			},
		}); err != nil {
			log.Fatalf("error requesting to join to DRLM Core: %v", err)
		}

		req, err := agentConn.Recv()
		if err != nil {
			log.Fatalf("error recieving the join response: %v", err)
		}

		switch req.MessageType {
		case drlm.AgentConnectionFromCore_MESSAGE_TYPE_JOIN_RESPONSE:
			if req.JoinResponse.Status != drlm.AgentConnectionFromCore_JoinResponse_STATUS_ACCEPT {
				log.Fatalf("core rejected the join request")
			}

			cfg.Config.Core.Secret = req.JoinResponse.CoreSecret
			cfg.Config.Minio.AccessKey = req.JoinResponse.MinioAccessKey
			cfg.Config.Minio.SecretKey = req.JoinResponse.MinioSecretKey

			log.Println("config saved!")

			// TODO: Save configuration!

		default:
			log.Fatalf("unknown core response when waiting for the join response: %s", req.MessageType.String())
		}

		agentConn, err = client.AgentConnection(prepareCtx(ctx))
		if err != nil {
			// TODO: RETRY
			log.Fatalf("error connecting to the core: %v", err)
		}
	}

	if err := agentConn.Send(&drlm.AgentConnectionFromAgent{
		MessageType: drlm.AgentConnectionFromAgent_MESSAGE_TYPE_CONN_ESTABLISH,
	}); err != nil {
		log.Fatalf("error establishing the connection with the core: %v", err)
	}

	c := make(chan models.JobUpdate)
	go func() {
		for {
			req, err := agentConn.Recv()
			if err != nil {
				// TODO: RETRY

				log.Fatalf("agent connection error: %v", err)
			}

			switch req.MessageType {
			case drlm.AgentConnectionFromCore_MESSAGE_TYPE_JOB_NEW:
				job.Run(c, req.JobNew.Id, req.JobNew.Name, req.JobNew.Config, req.JobNew.Target)

			case drlm.AgentConnectionFromCore_MESSAGE_TYPE_JOB_CANCEL:
				job.Cancel(req.JobCancel.Id)

			case drlm.AgentConnectionFromCore_MESSAGE_TYPE_INSTALL_BINARY:
				if err := binary.Install(req.InstallBinary.Bucket, req.InstallBinary.Name); err != nil {
					// TODO: Send a response back
					log.Errorf("error installing the binary: %v", err)
				}

			default:
				log.Errorf("unknown message type recieved from the DRLM Core: %s", req.MessageType.String())
			}
		}
	}()

	go func() {
		for {
			select {
			case u := <-c:
				agentConn.Send(&drlm.AgentConnectionFromAgent{
					MessageType: drlm.AgentConnectionFromAgent_MESSAGE_TYPE_JOB_UPDATE,
					JobUpdate: &drlm.AgentConnectionFromAgent_JobUpdate{
						JobId:  u.JobID,
						Status: u.Status,
						Info:   u.Info,
					},
				})
			}
		}
	}()

	log.Info("successfully connected to the DRLM Core")

	select {
	case <-ctx.Done():
		if conn != nil {
			// TODO: Close gracefully (cancel the jobs)
			agentConn.CloseSend()
			conn.Close()
		}
	}
}

func prepareCtx(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"api": API,
		"tkn": cfg.Config.Core.Secret,
	}))
}
