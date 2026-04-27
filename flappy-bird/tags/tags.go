package tags

import "github.com/yohamta/donburi"

const (
	TagPlayer     = "player"
	TagGround     = "ground"
	TagPipe       = "pipe"
	TagCollidable = "collidable"
)

var (
	Bird       = donburi.NewTag().SetName(TagPlayer)
	Ground     = donburi.NewTag().SetName(TagGround)
	Pipe       = donburi.NewTag().SetName(TagPipe)
	Collidable = donburi.NewTag().SetName(TagCollidable)
)
