package factory

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"

	"github.com/yohamta/donburi"
)

func CreatePlayer(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	player := archetype.Player.SpawnInWorld(w)

	width := float64(entity.Width)
	height := float64(entity.Height)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])

	r := resolv.NewRectangleFromCorners(X, Y, X+width, Y+height)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)
	component.Player.Set(player, &component.PlayerData{
		Speed:             8,
		SpeedAcceleration: 1.05,
		SpeedFriction:     0.94,
		Shape:             r,
	})

	sprite := project.GetSpriteByEntityInstance(entity)
	component.Sprite.SetValue(player, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return player
}
