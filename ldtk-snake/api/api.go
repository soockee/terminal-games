package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const (
	url = "https://stockhause.info:13337/score"
)

type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func GetScores() ([]Score, error) {
	// Define the URL of the server

	// Create an HTTP client
	client := &http.Client{}

	// Create a new HTTP request (GET in this example)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request:", err)
		return nil, err
	}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return nil, err
	}

	defer resp.Body.Close() // Close the response body after reading

	// Check the response status code (200 indicates success)
	if resp.StatusCode != http.StatusOK {
		slog.Error("Error:", slog.Int("code", resp.StatusCode))
		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body:", err)
		return nil, err
	}

	// Define a slice of Score structs to store the decoded data
	var scores []Score

	// Unmarshal the response body into the slice of Score structs
	err = json.Unmarshal(body, &scores)
	if err != nil {
		slog.Error("Error decoding JSON:", err)
		return nil, err
	}

	// Access and process the decoded data
	for _, score := range scores {
		slog.Info("Scoreinfo", slog.String("name", score.Name), slog.Int("Score", score.Score))
	}

	return scores, nil
}

func PostScore(score int, name string) error {
	// Create a Score struct with provided data
	scoreData := Score{Name: name, Score: score}

	// Marshal the Score struct into JSON bytes
	jsonData, err := json.Marshal(scoreData)
	if err != nil {
		return fmt.Errorf("error marshalling score data to JSON: %w", err)
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the API URL and JSON content type header
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending POST request: %w", err)
	}
	defer resp.Body.Close() // Close the response body after processing

	// Check the response status code (optional)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, string(body))
	}

	slog.Info("Score posted successfully!")
	return nil
}
