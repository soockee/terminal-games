package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
)

type Wall struct {
	entity *ldtkgo.Entity
	tile   *ebiten.Image
}

func NewWall() *Wall {
	return nil
}

func (w *Wall) Draw() {

}

func (w *Wall) Update() {

}
