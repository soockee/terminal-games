package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type BallData struct {
	Speed    float64
	Shape    *resolv.Circle
	MaxSpeed float64
}

var Ball = donburi.NewComponentType[BallData]()
