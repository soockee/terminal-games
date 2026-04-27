package tags

import "github.com/yohamta/donburi"

// Tag variables for entity queries.
var (
	Player     = donburi.NewTag().SetName("Player")
	Collidable = donburi.NewTag().SetName("Collidable")
	Ground     = donburi.NewTag().SetName("Ground")
)

// String constants for use in switch/label logic.
const (
	TagPlayer     = "Player"
	TagCollidable = "Collidable"
	TagGround     = "Ground"
)
