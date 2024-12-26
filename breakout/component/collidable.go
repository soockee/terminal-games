package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

// Collidable is something that can be collided with.
// Ask yourself, "Can I collide with this?"
// E.g. an ball can collide with a wall, but not with a button.
type CollidableData struct {
	Type  *donburi.ComponentType[donburi.Tag]
	Shape resolv.IShape
}

var Collidable = donburi.NewComponentType[CollidableData]()
