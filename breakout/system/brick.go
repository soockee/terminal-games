package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateBrick(ecs *ecs.ECS) {
	count := 0
	component.Brick.Each(ecs.World, func(e *donburi.Entry) {
		count++
	})
	if count == 0 {
		sceneData := component.SceneState.Get(component.SceneState.MustFirst(ecs.World))
		if next, ok := component.GetNextLevel(sceneData.CurrentScene); ok {
			event.SceneStateEvent.Publish(ecs.World, &event.SceneStateData{
				CurrentScene: component.LevelClearScene,
				NextScene:    next,
			})
		} else {
			event.SceneStateEvent.Publish(ecs.World, &event.SceneStateData{
				CurrentScene: component.GameOverScene,
			})
		}
	}
}

func DrawBrick(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Brick.Each(ecs.World, func(e *donburi.Entry) {
		b := component.Collidable.Get(e)

		component.DrawScaledSprite(screen, component.Sprite.Get(e).Images[0], b.Shape)
	})
}
