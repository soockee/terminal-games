package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/tags"
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

// NewPaddle creates a paddle entity from an LDtk entity.
func NewPaddle(w donburi.World, entity *ldtkgo.Entity, side component.PaddleSide, speed float64) *donburi.Entry {
	tlx, tly := entity.TopLeft()
	ew, eh := entity.Size()

	shape := resolv.NewRectangleFromTopLeft(float64(tlx), float64(tly), float64(ew), float64(eh))
	space := component.Space.Get(component.Space.MustFirst(w))
	space.Add(shape)

	e := w.Entry(w.Create(
		tags.Paddle,
		component.Paddle,
		component.Shape,
		component.Velocity,
		component.Sprite,
		component.Color,
	))

	component.Paddle.Set(e, &component.PaddleData{
		Speed: speed,
		Side:  side,
	})
	component.Shape.Set(e, &component.ShapeData{Shape: shape})
	component.Velocity.Set(e, &component.VelocityData{})
	component.Color.Set(e, &component.ColorData{Color: entity.ColorRGBA()})

	if sub := entity.SubImage(); sub != nil {
		component.Sprite.Set(e, &component.SpriteData{Image: ebiten.NewImageFromImage(sub)})
	}

	return e
}

// NewBall creates a ball entity from an LDtk entity using a circle collision shape.
func NewBall(w donburi.World, entity *ldtkgo.Entity, speed float64) *donburi.Entry {
	cx, cy := entity.Center()
	ew, eh := entity.Size()

	// Circle centered on the entity; radius = half the smaller dimension.
	radius := float64(min(ew, eh)) / 2
	shape := resolv.NewCircle(float64(cx), float64(cy), radius)
	space := component.Space.Get(component.Space.MustFirst(w))
	space.Add(shape)

	e := w.Entry(w.Create(
		tags.Ball,
		component.Ball,
		component.Shape,
		component.Velocity,
		component.Sprite,
		component.SpawnPos,
		component.Color,
	))

	component.Ball.Set(e, &component.BallData{Speed: speed, MaxSpeed: speed * 2})
	component.Shape.Set(e, &component.ShapeData{Shape: shape})
	component.Velocity.Set(e, &component.VelocityData{X: speed, Y: speed})
	component.SpawnPos.Set(e, &component.SpawnPosData{X: float64(cx), Y: float64(cy)})
	component.Color.Set(e, &component.ColorData{Color: entity.ColorRGBA()})

	if sub := entity.SubImage(); sub != nil {
		component.Sprite.Set(e, &component.SpriteData{Image: ebiten.NewImageFromImage(sub)})
	}

	return e
}

// NewWall creates a wall entity from an IntGrid cell rectangle.
// x, y, w, h are in pixels.
func NewWall(w donburi.World, x, y, width, height float64) *donburi.Entry {
	shape := resolv.NewRectangleFromTopLeft(x, y, width, height)
	space := component.Space.Get(component.Space.MustFirst(w))
	space.Add(shape)

	e := w.Entry(w.Create(
		tags.Wall,
		component.Shape,
	))

	component.Shape.Set(e, &component.ShapeData{Shape: shape})
	return e
}

// NewWallsFromIntGrid creates merged wall rectangles from IntGrid rows.
// Consecutive solid cells in a row are merged into a single wall entity.
func NewWallsFromIntGrid(w donburi.World, ig *ldtkgo.IntGrid) {
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
				NewWall(w, x, y, width, height)
				startCol = -1
			}
		}
	}
}
