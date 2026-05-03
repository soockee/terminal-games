package component

import (
	"github.com/soockee/terminal-games/super-mario-bros/physics"
	"github.com/yohamta/donburi"
)

// BodyData links an ECS entity to its physics body.
// Nil Body means the entity has no physics representation.
type BodyData struct {
	Body *physics.Body
}

var Body = donburi.NewComponentType[BodyData]()

// PhysicsSpaceData holds the physics.Space singleton. One per world.
type PhysicsSpaceData struct {
	Space *physics.Space
}

var PhysicsSpace = donburi.NewComponentType[PhysicsSpaceData]()
