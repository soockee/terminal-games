package archetype

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/tags"
	"github.com/yohamta/donburi"
)

// NewSpace creates a singleton entity holding the resolv.Space.
func NewSpace(w donburi.World, width, height, cellW, cellH int) {
	e := w.Entry(w.Create(component.Space))
	space := resolv.NewSpace(width, height, cellW, cellH)
	component.Space.Set(e, space)
}

// NewDebug creates the singleton debug toggle entity.
func NewDebug(w donburi.World) {
	e := w.Entry(w.Create(component.Debug))
	component.Debug.Set(e, &component.DebugData{Enabled: false})
}

// NewCamera creates the singleton camera entity with the given virtual screen scale.
func NewCamera(w donburi.World, scaleX, scaleY float64) {
	e := w.Entry(w.Create(component.Camera))
	component.Camera.Set(e, &component.CameraData{ScaleX: scaleX, ScaleY: scaleY})
}

// NewAudio creates the singleton audio entity with pre-decoded players.
func NewAudio(w donburi.World, ctx *audio.Context, bgm *audio.Player, sfx map[string]*audio.Player) {
	e := w.Entry(w.Create(component.Audio))
	component.Audio.Set(e, &component.AudioData{Ctx: ctx, BGMusic: bgm, SFX: sfx})
}

// SpawnEntities spawns all game entities from the LDtk level and world.
// It handles both level-placed entities (bird via tag) and definition-only
// entities (ground via project entity definition).
func SpawnEntities(w donburi.World, level *ldtkgo.Level, world *ldtkgo.World, cfg SpawnConfig) {
	// Singletons for game state.
	NewScore(w, cfg.PipeCount)
	NewGameOver(w)

	for _, ent := range level.EntitiesByTag(tags.TagPlayer) {
		NewBird(w, ent, cfg.BirdVelX, cfg.BirdJump, cfg.BirdGravity)
	}

	groundY := float64(level.Height)
	if def := world.Project.EntityDefByID("Ground"); def != nil {
		if sub := def.SubImage(world.Project); sub != nil {
			tileH := float64(def.Height)
			groundY = float64(level.Height) - tileH
			NewGround(w, ebiten.NewImageFromImage(sub), groundY, tileH)
		}
	}

	// Spawn pipe pairs from definition-only entity.
	if def := world.Project.EntityDefByID("Pipe"); def != nil {
		if sub := def.SubImage(world.Project); sub != nil {
			pipeImg := ebiten.NewImageFromImage(sub)
			x := float64(level.Width)
			for i := 0; i < cfg.PipeCount; i++ {
				if i > 0 {
					// random horizontal spacing between PipeMinGapHorizontal and 2x PipeInterval
					x += cfg.PipeMinGapHorizontal + rand.Float64()*(cfg.PipeInterval*2-cfg.PipeMinGapHorizontal)
				}
				// random vertical gap between PipeMinGapVertical and 2x PipeMinGapVertical
				minPipeH := 40.0
				gapSize := cfg.PipeMinGapVertical + rand.Float64()*cfg.PipeMinGapVertical
				halfGap := gapSize / 2
				minGapY := halfGap + minPipeH           // top pipe >= minPipeH
				maxGapY := groundY - halfGap - minPipeH // bottom pipe >= minPipeH
				if minGapY > maxGapY {
					minGapY = groundY / 2
					maxGapY = minGapY
				}
				gapCenterY := minGapY + rand.Float64()*(maxGapY-minGapY)
				NewPipePair(w, pipeImg, x, gapCenterY, gapSize, groundY)
			}
		}
	}
}

// NewBird creates the player bird entity from an LDtk entity.
func NewBird(w donburi.World, entity *ldtkgo.Entity, velX, jump, gravity float64) *donburi.Entry {
	tlx, tly := entity.TopLeft()
	ew, eh := entity.Size()
	cx, cy := entity.Center()

	shape := resolv.NewRectangleFromTopLeft(float64(tlx), float64(tly), float64(ew), float64(eh))
	space := component.Space.Get(component.Space.MustFirst(w))
	space.Add(shape)

	e := w.Entry(w.Create(
		tags.Bird,
		component.Shape,
		component.Sprite,
		component.SpawnPos,
		component.Player,
		component.Color,
	))

	component.Shape.Set(e, &component.ShapeData{Shape: shape})
	component.SpawnPos.Set(e, &component.SpawnPosData{X: float64(cx), Y: float64(cy)})
	component.Player.Set(e, &component.PlayerData{VelX: velX, Jump: jump, Gravity: gravity})
	component.Color.Set(e, &component.ColorData{Color: entity.ColorRGBA()})

	if sub := entity.SubImage(); sub != nil {
		component.Sprite.Set(e, &component.SpriteData{Image: ebiten.NewImageFromImage(sub)})
	}

	return e
}

// NewGround creates the ground singleton from a tile image and world Y position.
// The collision shape extends 1e7 px to the right — effectively infinite for the scroller.
func NewGround(w donburi.World, tile *ebiten.Image, groundY, tileH float64) *donburi.Entry {
	shape := resolv.NewRectangleFromTopLeft(0, groundY, 1e7, tileH)
	space := component.Space.Get(component.Space.MustFirst(w))
	space.Add(shape)

	e := w.Entry(w.Create(
		tags.Ground,
		tags.Collidable,
		component.Ground,
		component.Shape,
	))

	component.Ground.Set(e, &component.GroundData{Tile: tile, Y: groundY})
	component.Shape.Set(e, &component.ShapeData{Shape: shape})

	return e
}

// NewPipePair spawns a top and bottom pipe at world X with a gap centered at gapCenterY.
func NewPipePair(w donburi.World, pipeImg *ebiten.Image, x, gapCenterY, gapSize, groundY float64) {
	pipeW := float64(pipeImg.Bounds().Dx())

	// Bottom pipe: from gap bottom edge to ground.
	botTop := gapCenterY + gapSize/2
	botH := groundY - botTop
	if botH > 0 {
		newPipe(w, pipeImg, x, botTop, pipeW, botH)
	}

	// Top pipe: from screen top (y=0) down to gap top edge.
	topH := gapCenterY - gapSize/2
	if topH > 0 {
		e := newPipe(w, pipeImg, x, 0, pipeW, topH)
		component.Pipe.Get(e).FlipY = true
	}
}

func newPipe(w donburi.World, pipeImg *ebiten.Image, x, y, pw, ph float64) *donburi.Entry {
	shape := resolv.NewRectangleFromTopLeft(x, y, pw, ph)
	space := component.Space.Get(component.Space.MustFirst(w))
	space.Add(shape)

	e := w.Entry(w.Create(
		tags.Pipe,
		tags.Collidable,
		component.Shape,
		component.Sprite,
		component.Pipe,
	))

	component.Shape.Set(e, &component.ShapeData{Shape: shape})
	component.Sprite.Set(e, &component.SpriteData{Image: pipeImg})
	component.Pipe.Set(e, &component.PipeData{})

	return e
}

// NewScore creates the singleton score entity with a target for win detection.
func NewScore(w donburi.World, target int) {
	e := w.Entry(w.Create(component.Score))
	component.Score.Set(e, &component.ScoreData{Target: target})
}

// NewGameOver creates the singleton game-over state entity.
func NewGameOver(w donburi.World) {
	e := w.Entry(w.Create(component.GameState))
	component.GameState.Set(e, &component.GameOverData{})
}
