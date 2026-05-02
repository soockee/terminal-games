package component

import (
	"github.com/jakecoffman/cp/v2"
	"github.com/yohamta/donburi"
)

// BodyData links an ECS entity to a cp physics body and its shapes.
// W and H store the entity's pixel dimensions for coordinate conversion
// between cp (center-bottom) and ECS (top-left).
type BodyData struct {
	Body   *cp.Body
	Shapes []*cp.Shape
	W, H   float64
}

var Body = donburi.NewComponentType[BodyData]()

// PhysicsSpaceData holds the cp.Space singleton. One per world.
type PhysicsSpaceData struct {
	Space *cp.Space
}

var PhysicsSpace = donburi.NewComponentType[PhysicsSpaceData]()
