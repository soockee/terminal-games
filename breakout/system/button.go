package system

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateButton(ecs *ecs.ECS) {
	component.Button.Each(ecs.World, func(e *donburi.Entry) {
		// buttonObject := dresolv.GetObject(e)
		// buttonData := component.Button.Get(e)
	})
}

func HandleButtonClick(w donburi.World, e *event.Interaction) {
	switch e.Action {
	case component.ActionClick:
		component.Button.Each(w, func(entry *donburi.Entry) {
			b := component.Button.Get(entry)
			if isVecInObject(e.Position, b.Shape) {
				button := component.Button.Get(entry)
				button.HandlerFunc(w)
			}
		})
	}

}

func DrawButton(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Button.Each(ecs.World, func(e *donburi.Entry) {
		b := component.Button.Get(e)
		sprite := component.Sprite.Get(e)
		component.DrawScaledSprite(screen, sprite.Images[0], b.Shape)
	})
}

func isVecInObject(vec input.Vec, obj resolv.IShape) bool {
	hit := !obj.Bounds().Intersection(resolv.NewCircle(vec.X, vec.Y, 1).Bounds()).IsEmpty()
	if hit {
		slog.Debug("Hit", slog.Any("hit", hit))
	}
	return hit
}
