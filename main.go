package main

import (
	"context"
	"os"

	"github.com/cecobask/timescale-coding-challenge/cmd/root"
	"github.com/cecobask/timescale-coding-challenge/pkg/env"
	"github.com/cecobask/timescale-coding-challenge/pkg/log"
)

func main() {
	logger := log.DefaultLogger()
	logger.ExitOnError(env.LoadEnv())
	ctx := log.WithContext(context.Background(), logger)
	if err := root.NewCommand().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
