package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/pong/archetype"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const (
	ldtkDir     = "assets/ldtk"
	ballSpeed   = 4.0
	paddleSpeed = 5.0
)

// Game holds the ECS world and LDtk data.
type Game struct {
	ecs   *ecs.ECS
	world *ldtkgo.World

	screenW, screenH int

	// Pre-converted ebiten images for BG and layers.
	bgImage     *ebiten.Image
	layerImages []*ebiten.Image
	level       *ldtkgo.Level

	// Level switching keys 1-9.
	levelKeys []ebiten.Key
}

func main() {
	world, err := ldtkgo.LoadWorld("pong.ldtk", os.DirFS(ldtkDir))
	if err != nil {
		log.Fatalf("loading world: %v", err)
	}

	g := &Game{
		world: world,
		levelKeys: []ebiten.Key{
			ebiten.Key1, ebiten.Key2, ebiten.Key3,
			ebiten.Key4, ebiten.Key5, ebiten.Key6,
			ebiten.Key7, ebiten.Key8, ebiten.Key9,
		},
	}

	g.loadLevel(0)

	ebiten.SetWindowSize(g.screenW, g.screenH)
	ebiten.SetWindowTitle("Pong — ECS + ldtk-super-simple-loader")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// loadLevel creates a fresh ECS world from the given level index.
func (g *Game) loadLevel(index int) {
	if index < 0 || index >= len(g.world.Levels) {
		return
	}

	level := g.world.Levels[index]
	g.level = level
	g.screenW = level.Width
	g.screenH = level.Height

	// Fresh ECS world
	e := ecs.NewECS(donburi.NewWorld())

	// Register systems
	e.AddSystem(system.UpdateInput)
	e.AddSystem(system.UpdateMovement)
	e.AddSystem(system.UpdateCollision)
	e.AddSystem(system.UpdateScore)

	// Register renderers (layer 0 = default)
	e.AddRenderer(0, system.DrawEntities)
	e.AddRenderer(0, system.DrawScore)
	e.AddRenderer(0, system.DrawDebug)

	g.ecs = e

	// Create resolv space (cell size = 8 to match IntGrid granularity)
	archetype.NewSpace(e.World, level.Width, level.Height, 8, 8)
	archetype.NewDebug(e.World)

	// Spawn paddles from LDtk entities tagged "player"
	for _, ent := range level.EntitiesByTag("player") {
		side := component.SideLeft
		if ent.ID == "player_right" {
			side = component.SideRight
		}
		archetype.NewPaddle(e.World, ent, side, paddleSpeed)
	}

	// Spawn ball from LDtk entities tagged "ball"
	for _, ent := range level.EntitiesByTag("ball") {
		archetype.NewBall(e.World, ent, ballSpeed)
	}

	// Spawn wall entities from IntGrid
	if ig := level.IntGrid("collision_test"); ig != nil {
		spawnWallsFromIntGrid(e.World, ig)
	}

	// Prepare layer images
	g.bgImage = nil
	if level.BGImage != nil {
		g.bgImage = ebiten.NewImageFromImage(level.BGImage)
	}
	g.layerImages = nil
	for _, l := range level.LoadedLayers {
		g.layerImages = append(g.layerImages, ebiten.NewImageFromImage(l.Image))
	}
}

// spawnWallsFromIntGrid creates merged wall rectangles from IntGrid rows.
// It does a simple row-scan merge: consecutive solid cells in a row become one wall.
func spawnWallsFromIntGrid(w donburi.World, ig *ldtkgo.IntGrid) {
	gridSize := ig.Def.GridSize
	for row := 0; row < ig.Height; row++ {
		startCol := -1
		for col := 0; col <= ig.Width; col++ {
			solid := col < ig.Width && ig.At(col, row) != 0
			if solid && startCol < 0 {
				startCol = col
			}
			if !solid && startCol >= 0 {
				x := float64(startCol * gridSize)
				y := float64(row * gridSize)
				width := float64((col - startCol) * gridSize)
				height := float64(gridSize)
				archetype.NewWall(w, x, y, width, height)
				startCol = -1
			}
		}
	}
}

func (g *Game) Update() error {
	// Level switching: keys 1-9
	for i, key := range g.levelKeys {
		if i >= len(g.world.Levels) {
			break
		}
		if ebiten.IsKeyPressed(key) && g.world.Levels[i] != g.level {
			g.loadLevel(i)
			return nil
		}
	}

	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	if g.bgImage != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(g.bgImage.Bounds().Dx())
		bh := float64(g.bgImage.Bounds().Dy())
		op.GeoM.Scale(float64(g.screenW)/bw, float64(g.screenH)/bh)
		screen.DrawImage(g.bgImage, op)
	}

	// Layer images
	for _, img := range g.layerImages {
		screen.DrawImage(img, nil)
	}

	// ECS renderers (entities, score)
	g.ecs.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenW, g.screenH
}
