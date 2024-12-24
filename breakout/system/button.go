package system

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	dresolv "github.com/soockee/terminal-games/breakout/resolv"
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
		component.DrawScaledSprite(screen, component.Sprite.Get(e).Images[0], e)
	})
}

func isVecInObject(vec input.Vec, obj *resolv.ConvexPolygon) bool {
	// Y of object is skewed check objc creation
	slog.Debug("not yet implemented")

	return false
}
