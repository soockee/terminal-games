package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

type AnimationType int

type AnimationsData struct {
	Animation *ganim8.Animation
	Shape     resolv.IShape
	Loop    bool
}

var Animation = donburi.NewComponentType[AnimationsData]()
