package tags

import (
	"github.com/yohamta/donburi"
)

var (
	//GameElements
	Snake       = donburi.NewTag().SetName("Snake")
	SnakeBody   = donburi.NewTag().SetName("SnakeBody")
	Wall        = donburi.NewTag().SetName("Wall")
	Food        = donburi.NewTag().SetName("Food")
	Collidable  = donburi.NewTag().SetName("Collidable")
	Collectable = donburi.NewTag().SetName("Collectable")
	Mouse       = donburi.NewTag().SetName("Mouse")
	Animated    = donburi.NewTag().SetName("Animated")

	//UI
	Button    = donburi.NewTag().SetName("Button")
	TextField = donburi.NewTag().SetName("TextField")
)
