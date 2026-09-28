package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/yeixio/yggdrasil-core/internal/app"
	"github.com/yeixio/yggdrasil-core/internal/config"
	"github.com/yeixio/yggdrasil-core/internal/version"
)

func main() {
	dataDir := flag.String("data-dir", "", "application data directory")
	showVersion := flag.Bool("version", false, "print version, license, and corresponding source")
	flag.Parse()

	if *showVersion {
		fmt.Print(version.CurrentOffer().Text())
		os.Exit(0)
	}

	if *dataDir == "" {
		*dataDir = config.DefaultDataDir()
	}

	logger, logCloser, err := setupLogger(*dataDir)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	if logCloser != nil {
		defer logCloser.Close()
	}

	application, err := app.New(app.Options{DataDir: *dataDir, Logger: logger})
	if err != nil {
		logger.Error("failed to initialize", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Start(ctx); err != nil {
		logger.Error("daemon exited with error", "error", err)
		os.Exit(1)
	}
}

func setupLogger(dataDir string) (*slog.Logger, io.Closer, error) {
	logsDir := filepath.Join(dataDir, "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return nil, nil, err
	}
	logPath := filepath.Join(logsDir, "daemon.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, err
	}
	w := io.MultiWriter(os.Stdout, f)
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return logger, f, nil
}
