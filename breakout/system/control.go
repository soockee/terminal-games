package system

import (
	"slices"

	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/exp/maps"
)

func UpdateControl(ecs *ecs.ECS) {
	control := getControl(ecs)
	component.InputSytem.Update()

	// todo add follow cursor / touch press

	scene := component.SceneState.MustFirst(ecs.World)
	sceneState := component.SceneState.Get(scene)

	levels := maps.Keys(component.Levels)
	if slices.Contains(levels, sceneState.CurrentScene) {
		checkPlayerInput(ecs, control)
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

func checkPlayerInput(ecs *ecs.ECS, c *component.ControlData) {
	boost := false
	if c.InputHandler.ActionIsPressed(component.ActionMoveBoost) {
		boost = true
	}
	if ok := c.InputHandler.ActionIsPressed(component.ActionMoveLeft); ok {
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Action:    component.ActionMoveLeft,
			Direction: resolv.NewVector(-1, 0),
			Boost:     boost,
		})
	}
	if ok := c.InputHandler.ActionIsPressed(component.ActionMoveRight); ok {
		event.MoveEvent.Publish(ecs.World, &event.Move{
			Action:    component.ActionMoveRight,
			Direction: resolv.NewVector(1, 0),
			Boost:     boost,
		})
	}

}
