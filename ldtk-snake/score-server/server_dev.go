//go:build !prod
// +build !prod

package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (s *ApiServer) Run() {
	slog.Info("Running Dev Server")
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	serverLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}), slog.LevelDebug)

	router := http.DefaultServeMux
	router.HandleFunc("/", makeHTTPHandleFunc(s.handleIndex))
	router.HandleFunc("/score", makeHTTPHandleFunc(s.handleScore))

	loggingMiddleware := LoggingMiddleware(logger)
	loggedRouter := loggingMiddleware(router)

	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         "0.0.0.0:13337",
		Handler:      loggedRouter,
		ErrorLog:     serverLogger,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("Failed to start HTTP server", err)
		os.Exit(1)
	}
}
