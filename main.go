package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zaffka/aggregator/config"
	"github.com/zaffka/aggregator/exec"
	"github.com/zaffka/aggregator/log"
	"github.com/zaffka/aggregator/queue"
	"github.com/zaffka/aggregator/shutdown"
)

var (
	version    = "dev"
	configFile = ".configuration/configuration.json"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Init logger.
	logger := log.New(fmt.Sprintf("app-%s", version), os.Stderr)
	defer func() {
		err := logger.Sync()
		if err != nil {
			print(err)
		}
	}()

	// Getting configuration handler.
	configuration, err := config.Get(configFile)
	if err != nil {
		logger.Error(err)

		return
	}

	// Init new message queue.
	que := queue.NewQueue(
		ctx,
		configuration.QueueTopics(),
		configuration.QueueLen(),
		logger,
	)

	// Make new execution pool.
	pool := exec.NewPool(ctx, logger)

	// Assemble domain generators and aggregators from the configuration data.
	allAggrs := configuration.Aggregators()
	allGens := configuration.Generators()

	// Add to the execution pool every one of gens and aggrs using queue and a storage if needed.
	pool.Add(
		exec.For(allAggrs).WithQueue(que).WithStorage(nil).Build()...)

	pool.Add(
		exec.For(allGens).WithQueue(que).Build()...)

	// Lock till the interruption signal received.
	<-shutdown.Lock(cancelFn, 3*time.Second)

	logger.Info("exit")
}
