package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

type AnimationData struct {
	Animations map[int]*ganim8.Animation
}

var Animation = donburi.NewComponentType[AnimationData]()
