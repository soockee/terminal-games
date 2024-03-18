package game

import (
	"log/slog"
	"slices"

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
	"golang.org/x/exp/maps"
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
	g.start(component.StartScene, component.Empty)
	return g
}

func (g *Game) Update() error {
	if sceneState, ok := component.SceneState.First(g.ecs.World); ok {
		sceneStateData := component.SceneState.Get(sceneState)
		if g.scene.GetId() != sceneStateData.CurrentScene {
			g.reset()
			g.start(sceneStateData.CurrentScene, sceneStateData.LastScene)
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

func (g *Game) start(sceneId string, prevSceneId string) {

	// global systems
	g.ecs.AddSystem(system.UpdateControl)
	factory.CreateControl(g.ecs)
	factory.CreateSceneState(g.ecs, sceneId, prevSceneId, g.ldtkProject)
	factory.CreateSettings(g.ecs)

	g.ecs.AddRenderer(layers.Default, system.DrawDebug)
	g.ecs.AddRenderer(layers.Default, system.DrawHelp)

	pkgevents.SceneStateEvent.Subscribe(g.ecs.World, handleSceneStateEvent)
	pkgevents.GameStateEvent.Subscribe(g.ecs.World, handleGameStateEvent)

	g.scene = scene.CreateScene(sceneId, g.ecs, g.ldtkProject)
}

func (g *Game) reset() {
	// Check if sceneId is in SnakeLevels slice
	levels := maps.Keys(component.SnakeLevels)
	if slices.Contains(levels, g.scene.GetId()) {
		gamestate := component.GameState.MustFirst(g.ecs.World)
		gamedata := component.GameState.Get(gamestate)
		g.ecs = desc.NewECS(donburi.NewWorld())
		factory.FinalizeGameState(g.ecs, gamedata)
		return
	} else {
		g.ecs = desc.NewECS(donburi.NewWorld())
	}
}

func handleSceneStateEvent(w donburi.World, e *pkgevents.SceneStateData) {
	sceneStateData := component.SceneState.Get(component.SceneState.MustFirst(w))

	if e.CurrentScene == component.LevelClearScene {
		sceneStateData.LastScene = sceneStateData.CurrentScene
	}
	sceneStateData.CurrentScene = e.CurrentScene
}

func handleGameStateEvent(w donburi.World, e *pkgevents.GameStateData) {
	if e.IsGameOver {
		slog.Info("Game Over")
		sceneStateData := component.SceneState.Get(component.SceneState.MustFirst(w))
		sceneStateData.LastScene = sceneStateData.CurrentScene
		sceneStateData.CurrentScene = component.GameOverScene

	}
}
