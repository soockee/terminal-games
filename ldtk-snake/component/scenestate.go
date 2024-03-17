package component

import (
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/yohamta/donburi"
)

type SceneId int

const (
	Empty         SceneId = -1
	SnakeScene    SceneId = 0
	StartScene    SceneId = 1
	GameOverScene SceneId = 2
)

type SceneData struct {
	CurrentScene SceneId
	Project      *assets.LDtkProject
}

var SceneState = donburi.NewComponentType[SceneData]()
