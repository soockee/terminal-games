package component

import (
	"github.com/yohamta/donburi"
)

type SceneId int

const (
	Empty      SceneId = -1
	SnakeScene SceneId = 0
	StartScene SceneId = 1
)

type SceneDate struct {
	CurrentScene SceneId
}

var SceneState = donburi.NewComponentType[SceneDate]()
