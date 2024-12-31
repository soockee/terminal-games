package game

import (
	"log/slog"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	pkgevents "github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/scene"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/exp/maps"
)

var (
	store = &assets.LDtkProject{}
)

type Game struct {
	ecs   *ecs.ECS
	scene scene.Scene
	store *assets.LDtkProject
}

func NewGame(project *assets.LDtkProject) *Game {
	assets.MustLoadAssets()

	g := &Game{
		ecs:   ecs.NewECS(donburi.NewWorld()),
		store: project,
	}

	store = project

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
	level := g.store.Project.LevelByIdentifier(g.scene.GetId())
	opt := assets.NewDefaultDrawOptions()
	g.store.Renderer.Render(level, screen, opt)
	scene.Draw(g.scene, screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	w := g.store.Project.LevelByIdentifier(g.scene.GetId()).Width
	h := g.store.Project.LevelByIdentifier(g.scene.GetId()).Height
	return scene.Layout(w, h)
}

func (g *Game) start(sceneId string, prevSceneId string) {

	// global systems
	g.ecs.AddSystem(system.UpdateControl)
	g.ecs.AddSystem(system.UpdateAnimation)
	g.ecs.AddSystem(system.ProcessEvents)
	archetype.NewControl(g.ecs)
	archetype.NewSceneState(g.ecs, sceneId, prevSceneId, g.store)
	archetype.NewSettings(g.ecs)

	g.ecs.AddRenderer(layers.Default, system.DrawDebug)
	g.ecs.AddRenderer(layers.Default, system.DrawHelp)
	g.ecs.AddRenderer(layers.Default, system.DrawAnimation)

	pkgevents.SceneStateEvent.Subscribe(g.ecs.World, handleSceneStateEvent)
	pkgevents.GameStateEvent.Subscribe(g.ecs.World, handleGameStateEvent)
	pkgevents.CreateEntityEvent.Subscribe(g.ecs.World, handleCreateEntityEvent)

	g.scene = scene.CreateScene(sceneId, g.ecs, g.store)
}

func (g *Game) reset() {
	levels := maps.Keys(component.Levels)
	if slices.Contains(levels, g.scene.GetId()) {
		gamestate := component.GameState.MustFirst(g.ecs.World)
		gamedata := component.GameState.Get(gamestate)
		g.ecs = ecs.NewECS(donburi.NewWorld())
		archetype.AccumulateGameState(g.ecs, gamedata)
		return
	} else if g.scene.GetId() == component.LevelClearScene {
		gamestate := component.GameState.MustFirst(g.ecs.World)
		gamedata := component.GameState.Get(gamestate)
		g.ecs = ecs.NewECS(donburi.NewWorld())
		archetype.ContinueLevelGameState(g.ecs, gamedata)
		return
	} else {
		g.ecs = ecs.NewECS(donburi.NewWorld())
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

func handleCreateEntityEvent(w donburi.World, e *event.CreateEntityData) {
	if len(e.Tags) != 1 {
		slog.Error("Entity has not exactly one tag", slog.Any("tags", e.Tags))
		return
	}
	switch e.Tags[0] {
	case tags.Button.Name():
		shape := resolv.NewRectangleFromTopLeft(e.X, e.Y, e.W, e.H)
		sprite := store.GetSpriteByIdentifier(e.Identifier)
		archetype.NewButton(w, shape, sprite, e.Identifier)
	case tags.Player.Name():
		shape := resolv.NewRectangleFromTopLeft(e.X, e.Y, e.W, e.H)
		sprite := store.GetSpriteByIdentifier(e.Identifier)
		archetype.NewPlayer(w, shape, sprite)
	case tags.Ball.Name():
		shape := resolv.NewCircle(e.X, e.Y, e.W/2)
		sprite := ebiten.NewImage(int(shape.Bounds().Width()), int(shape.Bounds().Height()))
		archetype.NewBall(w, shape, sprite)
	case tags.Wall.Name():
		shape := resolv.NewRectangleFromTopLeft(e.X, e.Y, e.W, e.H)
		sprite := ebiten.NewImage(int(shape.Bounds().Width()), int(shape.Bounds().Height()))
		archetype.NewWall(w, shape, sprite)
	case tags.Brick.Name():
		shape := resolv.NewRectangleFromTopLeft(e.X, e.Y, e.W, e.H)
		sprite := store.GetSpriteByIdentifier(e.Identifier)
		archetype.NewBrick(w, shape, sprite)
	case tags.TextField.Name():
		shape := resolv.NewRectangle(e.X, e.Y, e.W, e.H)
		sprite := ebiten.NewImage(int(shape.Bounds().Width()), int(shape.Bounds().Height()))
		archetype.NewTextField(w, shape, sprite)
	case tags.Explosion.Name():
		shape := resolv.NewRectangle(e.X, e.Y, e.W, e.H)
		sprite, err := store.GetAnimatedSpriteByIdentifier(e.Identifier)
		if err != nil {
			slog.Debug("Failed animation data initialization", slog.String("entity", e.Identifier), slog.Any("error", err))
			panic(err)
		}
		archetype.NewExplosion(w, shape, *sprite)
	default:
		slog.Error("tag not found: noop", slog.Any("tag", e.Identifier))
	}
}
