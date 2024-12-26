package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type WallData struct {
	Shape *resolv.ConvexPolygon
}

var Wall = donburi.NewComponentType[WallData]()
