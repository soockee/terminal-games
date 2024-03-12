package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type VelocityData struct {
	Velocity resolv.Vector
}

var Velocity = donburi.NewComponentType[VelocityData]()
