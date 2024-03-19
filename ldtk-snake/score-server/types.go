package main

type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func NewScore(name string, score int) *Score {
	return &Score{
		Name:  name,
		Score: score,
	}
}
