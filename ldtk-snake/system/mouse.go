package system

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"
)

func UpdateMouse(ecs *ecs.ECS) {
	controlEntity := component.Control.MustFirst(ecs.World)
	controlData := component.Control.Get(controlEntity)

	query := donburi.NewQuery(filter.Contains(tags.Mouse))

	mouseEntry, ok := query.First(ecs.World)
	if !ok {
		slog.Error("could not find mouse in the world")
		panic(0)
	}
	mouseObj := component.Object.Get(mouseEntry)
	offsetX := mouseObj.Size.X / 2
	offsetY := mouseObj.Size.Y / 2
	var pos resolv.Vector
	if controlData.LastPosition != nil {
		pos = controlData.LastPosition.Sub(resolv.NewVector(offsetX, offsetY))
	} else {
		pos = resolv.Vector(controlData.InputHandler.CursorPos()).Sub(resolv.NewVector(offsetX, offsetY))
	}
	mouseObj.Position = pos

	mouseData := component.Mouse.Get(mouseEntry)
	mouseData.Invincible.Update()
}

func DrawMouse(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(tags.Mouse))

	mouseEntry, ok := query.First(ecs.World)
	if !ok {
		slog.Error("could not find mouse in the world")
		panic(0)
	}
	mouseObj := component.Object.Get(mouseEntry)
	mouseAnimation := component.Animation.Get(mouseEntry)
	mouseData := component.Mouse.Get(mouseEntry)

	offsetX := mouseObj.Size.X / 2
	offsetY := mouseObj.Size.Y / 2
	drawOptions := ganim8.DrawOpts(mouseObj.Position.X-offsetX, mouseObj.Position.Y-offsetY, 0, 2, 2)

	if mouseData.IsHidden {
		mouseAnimation.Animations[int(component.MouseHidden)].Draw(screen, drawOptions)
		mouseAnimation.Animations[int(component.MouseHidden)].Update()
	} else {
		mouseAnimation.Animations[int(component.MouseMoving)].Draw(screen, drawOptions)
		mouseAnimation.Animations[int(component.MouseMoving)].Update()
	}
}

func OnToggleMouse(w donburi.World, e *event.Mouse) {
	slog.Info("Toggle Mouse")
	entity := component.Mouse.MustFirst(w)
	mouseData := component.Mouse.Get(entity)
	mouseData.IsHidden = !mouseData.IsHidden
}
