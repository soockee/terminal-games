package main

import (
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

func main() {

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))

	slog.Info("Starting...")
	slog.Info("Setup Storage...")
	store, err := NewSQLiteStore()
	if err != nil {
		slog.Error("database initialization error", slog.Any("err", err))
	}
	server := NewApiServer("0.0.0.0:13337", store)
	server.Run()
}
