package tags

import (
	"github.com/yohamta/donburi"
)

var (
	//GameElements
	Player      = donburi.NewTag().SetName("Player")
	Wall        = donburi.NewTag().SetName("Wall")
	Collidable  = donburi.NewTag().SetName("Collidable")
	Collectable = donburi.NewTag().SetName("Collectable")
	Animated    = donburi.NewTag().SetName("Animated")

	//UI
	Button    = donburi.NewTag().SetName("Button")
	TextField = donburi.NewTag().SetName("TextField")
)