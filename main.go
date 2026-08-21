package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeonjw85/Ksecret/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := cmd.Execute(ctx); err != nil {
		if errors.Is(err, cmd.ErrHasFindings) {
			os.Exit(1)
		}
		slog.Error("ksecret", slog.Any("err", err))
		os.Exit(2)
	}
}
