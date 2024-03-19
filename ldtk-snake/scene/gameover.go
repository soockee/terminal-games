package scene

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	pkgevents "github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	"github.com/soockee/terminal-games/ldtk-snake/system"
	"github.com/soockee/terminal-games/ldtk-snake/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	decs "github.com/yohamta/donburi/ecs"
)

type GameOverScene struct {
	ecs         *decs.ECS
	ldtkProject *assets.LDtkProject
	once        *sync.Once
}

func NewGameOverScene(ecs *decs.ECS, project *assets.LDtkProject) *GameOverScene {
	return &GameOverScene{
		ecs:         ecs,
		ldtkProject: project,
		once:        &sync.Once{},
	}
}

func (s *GameOverScene) configure() {
	s.ecs.AddSystem(system.UpdateObjects)
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdateButton)

	s.ecs.AddRenderer(layers.Default, system.DrawDebug)
	s.ecs.AddRenderer(layers.Default, system.DrawHelp)
	s.ecs.AddRenderer(layers.Default, system.DrawButton)
	s.ecs.AddRenderer(layers.Default, system.DrawTextField)

	level := s.ldtkProject.Project.LevelByIdentifier(s.GetId())

	cellWidth := level.Width / level.Layers[layers.Default].CellWidth
	CellHeight := level.Height / level.Layers[layers.Default].CellHeight
	space := factory.CreateSpace(
		s.ecs,
		level.Width,
		level.Height,
		cellWidth,
		CellHeight,
	)

	CreateEntities(s, space)

	gamedata := component.GameState.Get(component.GameState.MustFirst(s.ecs.World))

	url := "http://stockhause.info:13337/score"
	scores, err := GetScores(url)
	highscore := util.CalculateHighscore(float64(gamedata.TotalScore), gamedata.TotalTime.Seconds())

	score := Score{
		Name:  util.GetRandomName(),
		Score: int(highscore),
	}
	scores = append(scores, score)

	slices.SortFunc(scores, func(a, b Score) int {
		if a.Score == b.Score {
			return 0
		} else if a.Score > b.Score {
			return -1
		}
		return 1
	})

	if err != nil {
		slog.Error("error fetching scores", err)
	}

	component.Text.Each(s.ecs.World, func(e *donburi.Entry) {
		textfield := component.Text.Get(e)
		if textfield.Identifier == "Score" {
			textfield.Text = append(textfield.Text, fmt.Sprintf("%d", highscore))
		} else if textfield.Identifier == "Time" {
			textfield.Text = append(textfield.Text, fmt.Sprintf("%.1fs", gamedata.TotalTime.Seconds()))
		} else if textfield.Identifier == "Leaderboard" {
			if len(scores) > 9 {
				scores = scores[:9]
			}
			for i, score := range scores {
				if len(score.Name) > 9 {
					score.Name = score.Name[:9]
				}
				textfield.Text = append(textfield.Text, fmt.Sprintf("%d. %s %d", i+1, score.Name, score.Score))
			}
		}
	})

	PostScore(score.Score, score.Name, url)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func GetScores(url string) ([]Score, error) {
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

func PostScore(score int, name string, apiURL string) error {
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
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
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

func (s *GameOverScene) GetId() string {
	return component.GameOverScene
}
func (s *GameOverScene) getLdtkProject() *assets.LDtkProject {
	return s.ldtkProject
}
func (s *GameOverScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *GameOverScene) getOnce() *sync.Once {
	return s.once
}
