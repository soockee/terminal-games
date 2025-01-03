package archetype

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/soockee/terminal-games/breakout/util"
	"github.com/yohamta/donburi"
)

var (
	Button = newArchetype(
		tags.Button,

		component.Sprite,
		component.Button,
	)

	buttonhandlerMapping = map[string]func(w donburi.World){
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
)

func NewButton(w donburi.World, shape resolv.IShape, sprite *ebiten.Image, buttonType string) *donburi.Entry {
	button := Button.SpawnInWorld(w)

	component.Space.Get(component.Space.MustFirst(w)).Add(shape)
	component.Button.Set(button, &component.ButtonData{
		Clicked:     false,
		HandlerFunc: buttonhandlerMapping[buttonType],
		Shape:       shape.(*resolv.ConvexPolygon),
		Type:        buttonType,
	})

	component.Sprite.SetValue(button, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return button
}
