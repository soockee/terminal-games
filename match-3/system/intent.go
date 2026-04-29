package system

// Intent represents a player action decoded from raw input.
type Intent interface {
	isIntent()
}

// SelectTile means the player wants to select/highlight a tile.
type SelectTile struct {
	Col, Row int
}

func (SelectTile) isIntent() {}

// InitiateSwap means the player wants to swap the selected tile with an adjacent one.
type InitiateSwap struct {
	FromCol, FromRow int
	ToCol, ToRow     int
}

func (InitiateSwap) isIntent() {}

// Deselect means the player tapped the already-selected tile.
type Deselect struct{}

func (Deselect) isIntent() {}

// ChangeSelection means the player tapped a non-adjacent tile while one is selected.
type ChangeSelection struct {
	Col, Row int
}

func (ChangeSelection) isIntent() {}

// StartGame means the player pressed action to start the game.
type StartGame struct{}

func (StartGame) isIntent() {}

// RestartGame means the player pressed action on win/loss screen.
type RestartGame struct{}

func (RestartGame) isIntent() {}
