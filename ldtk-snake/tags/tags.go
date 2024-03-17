package tags

import (
	"github.com/yohamta/donburi"
)

var (
	Snake      = donburi.NewTag().SetName("Snake")
	SnakeBody  = donburi.NewTag().SetName("SnakeBody")
	Wall       = donburi.NewTag().SetName("Wall")
	Food       = donburi.NewTag().SetName("Food")
	Button     = donburi.NewTag().SetName("Button")
	Collidable = donburi.NewTag().SetName("Collidable")
)
