package game

import (
	"log/slog"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	pkgevents "github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/factory"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/scene"
	"github.com/soockee/terminal-games/breakout/system"
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
	// Check if sceneId is in Levels slice
	levels := maps.Keys(component.Levels)
	if slices.Contains(levels, g.scene.GetId()) {
		gamestate := component.GameState.MustFirst(g.ecs.World)
		gamedata := component.GameState.Get(gamestate)
		g.ecs = desc.NewECS(donburi.NewWorld())
		factory.AccumulateGameState(g.ecs, gamedata)
		return
	} else if g.scene.GetId() == component.LevelClearScene {
		gamestate := component.GameState.MustFirst(g.ecs.World)
		gamedata := component.GameState.Get(gamestate)
		g.ecs = desc.NewECS(donburi.NewWorld())
		factory.ContinueLevelGameState(g.ecs, gamedata)
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
