package component

import (
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/yohamta/donburi"
)

type SceneId int

const (
	Empty                SceneId = -1
	SnakeScene           SceneId = 0
	SnakeBorderlessScene SceneId = 1
	StartScene           SceneId = 2
	GameOverScene        SceneId = 3
)

var SnakeLevels = []SceneId{
	SnakeScene,
	SnakeBorderlessScene,
}

type SceneData struct {
	CurrentScene SceneId
	Project      *assets.LDtkProject
}

var SceneState = donburi.NewComponentType[SceneData]()
