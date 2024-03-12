package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	pkgevents "github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	"github.com/soockee/terminal-games/ldtk-snake/scene"
	"github.com/soockee/terminal-games/ldtk-snake/system"
	"github.com/yohamta/donburi"
	desc "github.com/yohamta/donburi/ecs"
)

type Game struct {
	ecs         *desc.ECS
	scene       scene.Scene
	ldtkProject *assets.LDtkProject
}

func NewGame(project *assets.LDtkProject) *Game {
	assets.MustLoadAssets()

	g := &Game{
		ecs:         desc.NewECS(donburi.NewWorld()),
		ldtkProject: project,
	}
	g.start(component.StartScene)
	return g
}

func (g *Game) Update() error {
	if scenestate, ok := component.SceneState.First(g.ecs.World); ok {
		gamestateData := component.SceneState.Get(scenestate)
		if g.scene.GetId() != gamestateData.CurrentScene {
			g.reset()
			g.start(gamestateData.CurrentScene)
		}
	}
	scene.Update(g.scene)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	scene.Draw(g.scene, screen)

}

func (g *Game) Layout(width, height int) (int, int) {
	return scene.Layout(g.scene)
}

func (g *Game) start(sceneId component.SceneId) {

	// global systems
	g.ecs.AddSystem(system.UpdateControl)
	factory.CreateControl(g.ecs)
	factory.CreateSceneState(g.ecs, sceneId)
	factory.CreateSettings(g.ecs)

	g.ecs.AddRenderer(layers.Default, system.DrawDebug)
	g.ecs.AddRenderer(layers.Default, system.DrawHelp)

	pkgevents.SceneStateEvent.Subscribe(g.ecs.World, handleSceneStateEvent)

	g.scene = scene.CreateScene(sceneId, g.ecs, g.ldtkProject)
}

func (g *Game) reset() {
	settings, _ := component.Settings.First(g.ecs.World)
	settingsData := component.Settings.Get(settings)
	settingsData.Debug = false
	settingsData.ShowHelpText = false

	g.ecs = desc.NewECS(donburi.NewWorld())
}

func handleSceneStateEvent(w donburi.World, e *pkgevents.SceneStateData) {
	if scenestate, ok := component.SceneState.First(w); ok {
		gamestateData := component.SceneState.Get(scenestate)
		gamestateData.CurrentScene = e.CurrentScene
	}
}

func loadAssets() {
	//assets.LoadFont()
}
