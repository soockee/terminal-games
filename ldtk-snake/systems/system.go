package systems

import "github.com/hajimehoshi/ebiten/v2"

type System interface {
	Update()
	Draw(*ebiten.Image)
}
