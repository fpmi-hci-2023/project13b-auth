package main

import (
	"log/slog"
)

var (
	Version string
	Commit  string
)

func main() {
	slog.Info("auth service started with version: %s, commit: %s", Version, Commit)

	slog.Info("auth shutting down...")
}
