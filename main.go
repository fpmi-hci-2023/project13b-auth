package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	Version string
	Commit  string
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	sig, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("auth starting", slog.String("version", Version), slog.String("commit", Commit))

	go func() {
		<-sig.Done()
		stop()
	}()

	<-sig.Done()
	slog.Info("auth stopping")
}
