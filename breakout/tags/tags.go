package tags

import (
	"github.com/yohamta/donburi"
)

var (
	//GameElements
	Player      = donburi.NewTag().SetName("Player")
	Ball        = donburi.NewTag().SetName("Ball")
	Brick       = donburi.NewTag().SetName("Brick")
	Wall        = donburi.NewTag().SetName("Wall")
	Collidable  = donburi.NewTag().SetName("Collidable")
	Collectable = donburi.NewTag().SetName("Collectable")
	Animation   = donburi.NewTag().SetName("Animated")
	Explosion   = donburi.NewTag().SetName("Explosion")

	//UI
	Button    = donburi.NewTag().SetName("Button")
	TextField = donburi.NewTag().SetName("TextField")
)
