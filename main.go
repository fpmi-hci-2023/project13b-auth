package main

import (
	"context"
	"log/slog"
)

var (
	Version string
	Commit  string
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	slog.Info("auth service started with version: %s, commit: %s", Version, Commit)

	select {
	case <-ctx.Done():
		break
	}

	cancel()
	slog.Info("auth shutting down...")
}
