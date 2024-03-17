package system

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"

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
		component.Button.Each(w, func(entity *donburi.Entry) {
			buttonObject := dresolv.GetObject(entity)
			if isVecInObject(e.Position, buttonObject) {
				button := component.Button.Get(entity)
				button.HandlerFunc(w)
			}
		})
	}

}

func DrawButton(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Button.Each(ecs.World, func(e *donburi.Entry) {
		component.DrawScaledSprite(screen, e)
	})
}

func isVecInObject(vec input.Vec, obj *resolv.Object) bool {
	vecX, vecY := obj.Shape.Bounds()
	// Y of object is skewed check objc creation
	slog.Debug("Vecs", slog.Any("VecX", vecX), slog.Any("VecY", vecY), slog.Float64("w", vecY.X), slog.Float64("h", vecY.Y), slog.Any("click vec", vec))
	if vec.X >= vecX.X && vec.X <= vecY.X {
		if vec.Y >= vecX.Y && vec.Y <= vecY.Y {
			return true
		}
	}
	return false
}
