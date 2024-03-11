package factory

import (
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSceneState(ecs *ecs.ECS, scene component.SceneId) *donburi.Entry {
	scenestate := archetype.SceneState.Spawn(ecs)
	if scene != component.Empty {
		component.SceneState.SetValue(scenestate, component.SceneDate{
			CurrentScene: scene,
		})
	} else {
		component.SceneState.SetValue(scenestate, component.SceneDate{
			CurrentScene: component.Empty,
		})
	}
	return scenestate
}
