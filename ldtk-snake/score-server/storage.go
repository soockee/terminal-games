package main

import "database/sql"

type Storage interface {
	CreateScore(*Score) error
	GetScores() ([]*Score, error)
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore() (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", "./scores.db")
	if err != nil {
		return nil, err
	}

	createTable(db)

	return &SQLiteStore{
		db: db,
	}, nil
}

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS scores (
			id SERIAL PRIMARY KEY,
			score BIGINT,
			name VARCHAR(255)
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) CreateScore(score *Score) error {
	query := `
		INSERT INTO scores (score, name)
		VALUES ($1, $2);
	`
	_, err := s.db.Exec(query, score.Score, score.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStore) GetScores() ([]*Score, error) {
	query := `
		SELECT score, name
		FROM scores;
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := []*Score{}
	for rows.Next() {
		var score int
		var name string
		if err := rows.Scan(&score, &name); err != nil {
			return nil, err
		}
		scores = append(scores, NewScore(name, score))
	}
	return scores, nil
}
