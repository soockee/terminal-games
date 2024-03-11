package system

import (
	"log/slog"

	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/yohamta/donburi/ecs"
)

func UpdateControl(ecs *ecs.ECS) {
	control := getControl(ecs)
	component.InputSytem.Update()

	if control.InputHandler.ActionIsJustPressed(component.ActionMoveUp) {
		slog.Info("Publish Click")
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Direction: component.ActionMoveUp,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionMoveDown) {
		slog.Info("Publish Click")
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Direction: component.ActionMoveDown,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionMoveLeft) {
		slog.Info("Publish Click")
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Direction: component.ActionMoveLeft,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionMoveRight) {
		slog.Info("Publish Click")
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Direction: component.ActionMoveRight,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionClick) {
		slog.Info("Publish Click")
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Direction: component.ActionClick,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionDebug) {
		slog.Info("Publish Click")
		event.UpdateSettingEvent.Publish(ecs.World, &event.UpdateSetting{
			Action: component.ActionDebug,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionHelp) {
		slog.Info("Publish Click")
		event.UpdateSettingEvent.Publish(ecs.World, &event.UpdateSetting{
			Action: component.ActionHelp,
		})
	}
	if info, ok := control.InputHandler.JustPressedActionInfo(component.ActionClick); ok {
		slog.Info("Publish Click")
		event.InteractionEvent.Publish(ecs.World, &event.Interaction{
			Action:   component.ActionClick,
			Position: info.Pos,
		})
	}
}

func getControl(ecs *ecs.ECS) *component.ControlData {
	return component.Control.Get(
		component.Control.MustFirst(ecs.World),
	)
}
