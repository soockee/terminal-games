package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type ButtonData struct {
	Clicked     bool
	HandlerFunc func(donburi.World)
	Shape       *resolv.ConvexPolygon
}

var Button = donburi.NewComponentType[ButtonData]()
