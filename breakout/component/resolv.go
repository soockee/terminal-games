package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

var Circle = donburi.NewComponentType[resolv.Circle]()
var ConvexPolygon = donburi.NewComponentType[resolv.ConvexPolygon]()
var Space = donburi.NewComponentType[resolv.Space]()
