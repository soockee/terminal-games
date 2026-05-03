package component

import "github.com/yohamta/donburi"

// PlayerData holds player-specific movement state.
type PlayerData struct {
	GroundContacts int     // number of foot-sensor contacts with ground (>0 = grounded)
	MoveSpeed      float64 // horizontal speed in pixels/sec
	JumpForce      float64 // upward impulse magnitude
	MoveDir        float64 // -1 left, 0 idle, 1 right (set by input)
	JumpInput      bool    // true when jump was just pressed (set by input)
	DesiredVelX    float64 // horizontal velocity to inject during physics step
}

var Player = donburi.NewComponentType[PlayerData]()
