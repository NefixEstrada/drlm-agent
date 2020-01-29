// SPDX-License-Identifier: AGPL-3.0-only

package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/brainupdaters/drlm-agent/transport/grpc"

	log "github.com/sirupsen/logrus"
)

// Main is the main function of DRLM Agent
func Main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "wg", &wg)

	go grpc.Init(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	select {
	case <-stop:
		fmt.Println("")
		log.Info("stopping DRLM Agent...")

		cancel()
		wg.Wait()
	}
}
