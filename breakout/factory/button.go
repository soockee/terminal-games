package factory

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/util"
	"github.com/yohamta/donburi"
)

var buttonhandlerMapping = map[string]func(w donburi.World){
	"StartButton": func(w donburi.World) {
		event.SceneStateEvent.Publish(w, &event.SceneStateData{
			CurrentScene: component.Level_0,
		})
	},
	"GithubButton": func(w donburi.World) { util.OpenUrl("https://github.com/soockee") },
	"ResetButton": func(w donburi.World) {
		for k := range component.Levels {
			component.Levels[k] = false
		}
		event.SceneStateEvent.Publish(w, &event.SceneStateData{
			CurrentScene: component.Level_0,
		})
	},
	"NextButton": func(w donburi.World) {
		sceneentry := component.SceneState.MustFirst(w)
		scenedata := component.SceneState.Get(sceneentry)
		if next, ok := component.GetNextLevel(scenedata.LastScene); ok {
			slog.Info("starting level", slog.Any("Level", next))
			event.SceneStateEvent.Publish(w, &event.SceneStateData{
				CurrentScene: next,
			})
		}
	},
}

func CreateButton(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	button := archetype.Button.SpawnInWorld(w)

	width := float64(entity.Width)
	height := float64(entity.Height)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])

	r := resolv.NewRectangleFromCorners(X, Y, X+width, Y+height)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)
	component.Button.Set(button, &component.ButtonData{
		Clicked:     false,
		HandlerFunc: buttonhandlerMapping[entity.Identifier],
		Shape:       r,
	})

	sprite := project.GetSpriteByEntityInstance(entity)

	component.Sprite.SetValue(button, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return button
}
