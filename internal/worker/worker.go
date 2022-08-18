package worker

import (
	"context"

	"github.com/shipt/tempest/logging"
)

// Config ...
type Config struct{}

// Run ...
func Run(ctx context.Context, cfg Config) error {
	logger := logging.New()
	logger.Info(ctx, "starting worker")

	// block until the parent context is cancelled
	// in reality, you'd likely kick off some kind of
	// message consumer here ...
	<-ctx.Done()
	return nil
}
