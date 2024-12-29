package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/ganim8/v2"
)

type AnimationType int

const (
	BrickExplosion AnimationType = iota
)

type Animation struct {
	Animation *ganim8.Animation
	Position  math.Vec2
}

type AnimationsData map[AnimationType]Animation

var Animations = donburi.NewComponentType[AnimationsData]()
