package archetype

import (
	"github.com/soockee/terminal-games/super-mario-bros/assets"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
)

// NewPlayer creates the player ECS entity.
func NewPlayer(w donburi.World, x, y, width, height float64) {
	entry := w.Entry(w.Create(
		component.Position,
		component.Animation,
		component.Player,
		component.Body,
	))
	component.Position.Set(entry, &component.PositionData{X: x, Y: y})
	anim := assets.PlayerAnimations()
	component.Animation.Set(entry, &anim)
	component.Player.Set(entry, &component.PlayerData{
		MoveSpeed: 120,
		JumpForce: 350,
	})

	spaceEntry, ok := component.PhysicsSpace.First(w)
	if !ok {
		return
	}
	space := component.PhysicsSpace.Get(spaceEntry).Space
	pd := component.Player.Get(entry)
	body := space.AddPlayerBody(x, y, width, height, &pd.DesiredVelX)
	component.Body.Set(entry, &component.BodyData{Body: body})
}
