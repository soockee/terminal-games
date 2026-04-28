package archetype

// SpawnConfig holds the game-design tuning knobs for entity creation.
type SpawnConfig struct {
	BirdVelX    float64 // horizontal bird velocity
	BirdJump    float64 // upward impulse on jump
	BirdGravity float64 // downward acceleration per tick

	PipeMinGapVertical   float64 // minimum vertical gap between top and bottom pipe (pixels)
	PipeMinGapHorizontal float64 // minimum horizontal gap between pipes (pixels)
	PipeInterval         float64 // horizontal distance between pipe pairs (pixels)
	PipeCount            int     // number of pipe pairs to spawn ahead
}
