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
		component.DrawScaledSprite(screen, component.Sprite.Get(e).Images[0], b.Shape)
	})
}

func isVecInObject(vec input.Vec, obj resolv.IShape) bool {
	objVec := obj.Bounds()
	// Y of object is skewed check objc creation
	slog.Debug("Vecs", slog.Any("VecX", objVec.Width()), slog.Any("VecY", objVec.Height()), slog.Float64("w", objVec.Width()), slog.Float64("h", objVec.Height()), slog.Any("click vec", vec))

	hit := !obj.Bounds().Intersection(resolv.NewCircle(vec.X, vec.Y, 1).Bounds()).IsEmpty()
	if hit {
		slog.Debug("Hit", slog.Any("hit", hit))
	}
	return hit
}
