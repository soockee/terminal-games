package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	Speed             float64
	SpeedFriction     float64
	SpeedAcceleration float64
	Shape             *resolv.ConvexPolygon
}

var Player = donburi.NewComponentType[PlayerData]()
