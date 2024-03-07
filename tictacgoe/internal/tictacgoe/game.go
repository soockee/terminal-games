package tictacgoe

import (
	"bytes"

	"github/soockee/terminal-games/tictacgoe/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/font"

	"github.com/eihigh/canvas"
	renderer "github.com/eihigh/canvas/renderers/ebiten"
)

// Game represents a game state.
type Game struct {
	input                *Input
	board                *Board
	turn                 *Turn
	boardImage           *ebiten.Image
	winImage             *ebiten.Image
	turnBoxImage         *ebiten.Image
	resetButtonImage     *ebiten.Image
	performanceInfoImage *ebiten.Image
	ScreenWidth          int
	ScreenHeight         int
	BoardSize            int
	gameState            *Gamestate
	gameOver             bool
	win                  *Win
	resetButton          *ResetButton
	performanceInfo      *PerformanceInfo
	font                 font.Face
}

var (
	tileImage   = ebiten.NewImage(tileSize, tileSize)
	crossImage  = ebiten.NewImage(tileSize, tileSize)
	circleImage = ebiten.NewImage(tileSize, tileSize)
	evenImage   = ebiten.NewImage(tileSize, tileSize)
)

func init() {
	var err error
	crossImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(assets.OnionBoy_png))
	if err != nil {
		log.Fatal().AnErr("images could not be loaded", err)
	}
	circleImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(assets.Pig_png))
	if err != nil {
		log.Fatal().AnErr("images could not be loaded", err)
	}
	evenImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(assets.Handshake_gif))
	if err != nil {
		log.Fatal().AnErr("images could not be loaded", err)
	}
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	config := NewConfig()
	g := &Game{
		input:           NewInput(),
		ScreenWidth:     config.screenWidth,
		ScreenHeight:    config.screenHeight,
		BoardSize:       config.boardSize,
		gameState:       NewGamestate(),
		gameOver:        false,
		win:             NewWin(),
		performanceInfo: NewPerformanceInfo(),
		font:            config.font,
	}
	var err error
	g.board, err = NewBoard(config.boardSize)
	if err != nil {
		return nil, err
	}
	g.turn = NewTurn(g.gameState.getCurrentPlayer())
	g.resetButton = NewResetButton(*g.input)
	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.input.Update()
	if g.gameOver {
		g.resetButton.Update(g)
		return nil
	}
	if g.gameState.moveCounter == 9 {
		g.gameOver = true
		return nil
	}
	// check tile clicked
	if g.input.mouseState == mouseStateSettled {
		offsetX, offsetY := g.GetBoardTranslation()
		if 0 > g.input.mouseInitPosX-offsetX || 0 > g.input.mouseInitPosY-offsetY {
			return nil
		}
		bw, bh := g.board.Size()
		row := CalculateCell(g.input.mouseInitPosX-offsetX, bw, g.board.size)
		column := CalculateCell(g.input.mouseInitPosY-offsetY, bh, g.board.size)

		move := column*g.BoardSize + row

		if move > 8 || move < 0 {
			return nil
		}

		if ((g.gameState.boards[(g.gameState.getCurrentPlayer()^1)] & bitfilterForBoard) & (0b01 << move)) == 0 {
			g.board.tiles[move].current.TileState = TileState(g.gameState.getCurrentPlayer())
			g.gameState.MakeMove(move)
			g.turn.turnState = g.gameState.getCurrentPlayer()
			if g.gameState.IsWin() && !g.gameOver {
				g.gameOver = true
				return nil
			}
		}

	}
	return nil
}

func CalculateCell(relativeMouseClickPosition int, boardLength int, cells int) int {
	d := boardLength - (boardLength - relativeMouseClickPosition)
	return (d / tileSize)
}

func (g *Game) GetBoardTranslation() (int, int) {
	boardWidth, boardHeight := g.boardImage.Size()
	x := (g.ScreenWidth - boardWidth) / 2
	y := (g.ScreenHeight - boardHeight) / 2
	return x, y
}

func (g *Game) GetResetButtonTranslation() (int, int) {
	resetButtonWidth, _ := g.resetButton.Size()
	x := (g.ScreenWidth - resetButtonWidth) / 2
	y := (g.ScreenHeight - ResetButtonBottomMargin)
	return x, y
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	// check win condition and end game with new screen if so
	if g.gameOver {
		if !g.gameState.IsWin() {
			g.DrawWin(screen, evenImage, -1)
		} else {
			switch g.gameState.getCurrentPlayer() {
			case 0:
				g.DrawWin(screen, crossImage, 0)
			case 1:
				g.DrawWin(screen, circleImage, 1)
			}
		}
		g.DrawBoard(screen)
		g.DrawTurnBox(screen)
		g.DrawResetButton(screen)
	} else {
		g.DrawBoard(screen)
		g.DrawTurnBox(screen)
	}

	g.DrawPerformanceInfo(screen)
}

func (g *Game) DrawBoard(screen *ebiten.Image) {
	if g.boardImage == nil {
		w, h := g.board.Size()
		g.boardImage = ebiten.NewImage(w, h)
	}
	op := &ebiten.DrawImageOptions{}
	x, y := g.GetBoardTranslation()
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)
	g.board.Draw(g.boardImage)
}

func (g *Game) DrawWin(screen *ebiten.Image, winnerImage *ebiten.Image, winner int) {
	if g.winImage == nil {
		w, h := g.win.Size()
		g.winImage = ebiten.NewImage(w, h)
	}
	op := &ebiten.DrawImageOptions{}
	winWidth, _ := g.win.Size()
	x := (g.ScreenWidth - winWidth) / 2
	y := (WinBoxTopMargin)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.winImage, op)
	g.win.Draw(g.winImage, winnerImage, winner, g.font)
}

func (g *Game) DrawTurnBox(screen *ebiten.Image) {
	if g.turnBoxImage == nil {
		width, height := g.turn.Size()
		g.turnBoxImage = ebiten.NewImage(width, height)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(TurnBoxMargin), float64(g.ScreenHeight-TurnBoxMargin))
	op.ColorM.ScaleWithColor(frameColor)
	screen.DrawImage(g.turnBoxImage, op)
	g.turn.Draw(g.turnBoxImage)
}

func (g *Game) DrawResetButton(screen *ebiten.Image) {
	if g.resetButtonImage == nil {
		width, height := g.resetButton.Size()
		g.resetButtonImage = ebiten.NewImage(width, height)
	}

	r := renderer.New(g.resetButtonImage)
	ctx := canvas.NewContext(r)
	x, y := g.GetResetButtonTranslation()
	DrawCanvasEclipse(ctx)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.resetButtonImage, op)
	g.resetButton.Draw(g.resetButtonImage, g.font)

}

func (g *Game) DrawPerformanceInfo(screen *ebiten.Image) {
	if g.performanceInfoImage == nil {
		width, height := g.performanceInfo.Size()
		g.performanceInfoImage = ebiten.NewImage(width, height)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.performanceInfo.Margin()), float64(g.performanceInfo.Margin()))
	op.ColorM.ScaleWithColor(frameColor)
	screen.DrawImage(g.performanceInfoImage, op)
	g.performanceInfo.Draw(g.performanceInfoImage, g.font)
}

func DrawCanvasEclipse(c *canvas.Context) {
	// Draw an closed set of points being smoothed\
	polyline := &canvas.Polyline{}
	polyline.Add(c.Width()/2, 0.0)
	polyline.Add(c.Width(), -1*c.Height()/2)
	polyline.Add(c.Width()/2, -1*c.Height())
	polyline.Add(0.0, -1*c.Height()/2)
	polyline.Add(c.Width()/2, 0.0)
	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(1)
	c.DrawPath(0, 0, polyline.Smoothen().Scale(1, -1))
	c.SetStrokeColor(frameColor)
	c.SetStrokeWidth(3)
	c.DrawPath(0, 0, polyline.Smoothen().Scale(1, -1))
}

func GetMessageDrawLength(msg string, font font.Face) int {
	width := 0
	for _, letter := range msg {
		bounds, _, ok := font.GlyphBounds(letter)
		if !ok {
			log.Fatal().Msg("glype not found")
		}
		width += int(bounds.Max.X.Round())
	}
	return width
}

func GetMessageMaxHeight(msg string, font font.Face) int {
	maxHeight := 0
	for _, letter := range msg {
		bounds, _, ok := font.GlyphBounds(letter)
		if !ok {
			log.Fatal().Msg("glype not found")
		}
		if maxHeight < bounds.Max.Y.Round() {
			maxHeight = bounds.Max.Y.Round()
		}
	}
	return maxHeight
}
