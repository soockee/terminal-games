package game

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	pkgevents "github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/scene"
	"github.com/soockee/terminal-games/ldtk-snake/system"
	"github.com/yohamta/donburi"
	desc "github.com/yohamta/donburi/ecs"
)

type Scene interface {
	ebiten.Game
	GetId() component.Scene
}
type Game struct {
	ecs   *desc.ECS
	scene Scene
}

func NewGame() *Game {
	g := &Game{
		ecs: desc.NewECS(donburi.NewWorld()),
	}
	g.scene = scene.NewStartScene(g.ecs)
	g.start()

	return g
}

func (g *Game) Update() error {
	g.scene.Update()
	if gamestate, ok := component.Gamestate.First(g.ecs.World); ok {
		gamestateData := component.Gamestate.Get(gamestate)
		if gamestateData.CurrentScene == g.scene.GetId() {
			return nil
		}
		if gamestateData.CurrentScene == component.SnakeScene {
			g.reset()
			g.scene = scene.NewSnakeScene(g.ecs)
		} else if gamestateData.CurrentScene == component.StartScreen {
			g.reset()
			g.scene = scene.NewStartScene(g.ecs)
		} else {
			slog.Error("invalid game state")
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	return config.C.LDtkProject.Levels[config.C.CurrentLevel].Width, config.C.LDtkProject.Levels[config.C.CurrentLevel].Height
}

func (g *Game) start() {
	g.ecs.AddSystem(system.UpdateControl)
	factory.CreateControl(g.ecs)
	factory.CreateGamestate(g.ecs, nil)

	pkgevents.GamestateEvent.Subscribe(g.ecs.World, handleGamestateEvent)
}

func (g *Game) reset() {
	g.ecs = desc.NewECS(g.ecs.World)
	g.start()
}

func handleGamestateEvent(w donburi.World, e *event.Gamestate) {
	if gamestate, ok := component.Gamestate.First(w); ok {
		gamestateData := component.Gamestate.Get(gamestate)
		gamestateData.CurrentScene = e.CurrentScene
	}
}
