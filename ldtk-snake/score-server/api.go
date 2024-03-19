package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

type ApiServer struct {
	store  Storage
	ipaddr string
}

func NewApiServer(ipaddr string, store Storage) *ApiServer {
	server := &ApiServer{
		store:  store,
		ipaddr: ipaddr,
	}
	return server
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func (s *ApiServer) Run() {
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
		Addr:         s.ipaddr,
		Handler:      loggedRouter,
		ErrorLog:     serverLogger,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("Failed to start HTTP server", err)
		os.Exit(1)
	}
}

func (s *ApiServer) handleIndex(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		WriteJson(w, http.StatusOK, "This is the score server for snake by soockee")
	default:
		return errors.New("method not allowed")
	}
	return nil
}

func (s *ApiServer) handleScore(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		scores, err := s.store.GetScores()
		if err != nil {
			return err
		}
		WriteJson(w, http.StatusOK, scores)
		return nil
	case "POST":
		score := &Score{}
		if err := json.NewDecoder(r.Body).Decode(score); err != nil {
			return err
		}
		s.store.CreateScore(score)
		WriteJson(w, http.StatusCreated, score)
	default:
		return errors.New("method not allowed")
	}
	return nil
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return cors(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	})
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		next(w, r)
	}
}
