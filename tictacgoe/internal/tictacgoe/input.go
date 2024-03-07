package tictacgoe

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type mouseState int
type spaceState int

const (
	mouseStateNone     mouseState = 0
	mouseStatePressing mouseState = 1
	mouseStateSettled  mouseState = 2

	spaceStateNone     spaceState = 0
	spaceStatePressing spaceState = 1
	spaceStateSettled  spaceState = 2
)

// Input represents the current key states.
type Input struct {
	mouseState    mouseState
	mouseInitPosX int
	mouseInitPosY int

	spaceState spaceState
}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{}
}

// Update updates the current input states.
func (i *Input) Update() {
	switch i.mouseState {
	case mouseStateNone:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.mouseInitPosX = x
			i.mouseInitPosY = y
			i.mouseState = mouseStatePressing
		}
	case mouseStatePressing:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			i.mouseState = mouseStateSettled
		}
	case mouseStateSettled:
		i.mouseState = mouseStateNone
	}

	switch i.spaceState {
	case spaceStateNone:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			i.spaceState = spaceStatePressing
		}
	case spaceStatePressing:
		if !ebiten.IsKeyPressed(ebiten.KeySpace) {
			i.spaceState = spaceStateSettled
		}
	case spaceStateSettled:
		i.spaceState = spaceStateNone
	}
}

func (i *Input) ToString() string {
	return fmt.Sprintf("mouseState: %v mouseInitPosX: %v mouseInitPosY: %v", i.mouseState, i.mouseInitPosX, i.mouseInitPosY)
}
