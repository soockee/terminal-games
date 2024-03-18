package system

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/yohamta/donburi/ecs"
)

func UpdateControl(ecs *ecs.ECS) {
	control := getControl(ecs)
	component.InputSytem.Update()

	if control.InputHandler.ActionIsJustPressed(component.ActionMoveUp) {
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Action: component.ActionMoveUp,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionMoveDown) {
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Action: component.ActionMoveDown,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionMoveLeft) {
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Action: component.ActionMoveLeft,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionMoveRight) {
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Action: component.ActionMoveRight,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionMoveHalt) {
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Action: component.ActionMoveHalt,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionDebug) {
		event.UpdateSettingEvent.Publish(ecs.World, &event.UpdateSetting{
			Action: component.ActionDebug,
		})
	}
	if control.InputHandler.ActionIsJustPressed(component.ActionHelp) {
		event.UpdateSettingEvent.Publish(ecs.World, &event.UpdateSetting{
			Action: component.ActionHelp,
		})
	}
	if info, ok := control.InputHandler.JustPressedActionInfo(component.ActionClick); ok {
		event.InteractionEvent.Publish(ecs.World, &event.Interaction{
			Action:   component.ActionClick,
			Position: info.Pos,
		})
	}
}

func getControl(ecs *ecs.ECS) *component.ControlData {
	return component.Control.Get(component.Control.MustFirst(ecs.World))
}
