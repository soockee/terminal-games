package system

import (
	"log/slog"
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
			if info, ok := control.InputHandler.JustPressedActionInfo(component.ActionClick); ok {
				control.LastPosition = (*resolv.Vector)(&info.Pos)
			}
		} else {
			cursorPosition := control.InputHandler.CursorPos()
			boost := false
			if control.InputHandler.ActionIsPressed(component.ActionMoveBoost) {
				boost = true
			}
			slog.Info("Cursor Info", slog.Float64("now X", cursorPosition.X), slog.Float64("now Y", cursorPosition.Y), slog.Float64("last X", control.LastPosition.X), slog.Float64("last Y", control.LastPosition.Y))
			if control.InputHandler.ActionIsPressed(component.ActionClick) {
				control.LastPosition = (*resolv.Vector)(&cursorPosition)
				event.MoveEvent.Publish(ecs.World, &event.Move{
					Action:   component.ActionMovePosition,
					Position: resolv.NewVector(cursorPosition.X, cursorPosition.Y),
					Boost:    boost,
				})
			}
		}
	}

	// if control.InputHandler.ActionIsJustPressed(component.ActionMoveUp) {
	// 	event.MoveEvent.Publish(ecs.World, &event.Move{
	// 		Action: component.ActionMoveUp,
	// 	})
	// }
	// if control.InputHandler.ActionIsJustPressed(component.ActionMoveDown) {
	// 	event.MoveEvent.Publish(ecs.World, &event.Move{
	// 		Action: component.ActionMoveDown,
	// 	})
	// }
	// if control.InputHandler.ActionIsJustPressed(component.ActionMoveLeft) {
	// 	event.MoveEvent.Publish(ecs.World, &event.Move{
	// 		Action: component.ActionMoveLeft,
	// 	})
	// }
	// if control.InputHandler.ActionIsJustPressed(component.ActionMoveRight) {
	// 	event.MoveEvent.Publish(ecs.World, &event.Move{
	// 		Action: component.ActionMoveRight,
	// 	})
	// }
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
