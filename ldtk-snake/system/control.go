package system

import (
	"slices"

	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/exp/maps"
)

func UpdateControl(ecs *ecs.ECS) {
	control := getControl(ecs)
	component.InputSytem.Update()

	// todo add follow cursor / touch press

	scene := component.SceneState.MustFirst(ecs.World)
	sceneState := component.SceneState.Get(scene)

	levels := maps.Keys(component.SnakeLevels)
	if slices.Contains(levels, sceneState.CurrentScene) {
		if control.LastPosition == nil {
			if info, ok := control.InputHandler.JustPressedActionInfo(component.ActionMovePosition); ok {
				control.LastPosition = (*resolv.Vector)(&info.Pos)
			}
		} else {
			boost := false
			if control.InputHandler.ActionIsPressed(component.ActionMoveBoost) {
				boost = true
			}
			// toggle mouse state on press and release
			if control.InputHandler.ActionIsJustPressed(component.ActionMovePosition) {
				event.MouseEvent.Publish(ecs.World, &event.Mouse{})
			}
			if control.InputHandler.ActionIsJustReleased(component.ActionMovePosition) {
				event.MouseEvent.Publish(ecs.World, &event.Mouse{})
			}
			// check continuously for postition
			if info, ok := control.InputHandler.PressedActionInfo(component.ActionMovePosition); ok {
				control.LastPosition = (*resolv.Vector)(&info.Pos)
				event.MoveEvent.Publish(ecs.World, &event.Move{
					Action:   component.ActionMovePosition,
					Position: resolv.NewVector(info.Pos.X, info.Pos.Y),
					Boost:    boost,
				})
			}
		}
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
