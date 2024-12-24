package factory

import (
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSceneState(ecs *ecs.ECS, scene string, lastScene string, project *assets.LDtkProject) *donburi.Entry {
	scenestate := archetype.SceneState.Spawn(ecs)
	if scene != component.Empty {
		component.SceneState.SetValue(scenestate, component.SceneData{
			LastScene:    lastScene,
			CurrentScene: scene,
			Project:      project,
		})
	} else {
		component.SceneState.SetValue(scenestate, component.SceneData{
			LastScene:    lastScene,
			CurrentScene: component.Empty,
			Project:      project,
		})
	}
	return scenestate
}
