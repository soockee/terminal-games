package component

import (
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/yohamta/donburi"
)

type SceneId int

const (
	Empty         SceneId = -1
	Level_0       SceneId = 0
	Level_1       SceneId = 1
	Level_2       SceneId = 2
	Level_3       SceneId = 3
	StartScene    SceneId = 4
	GameOverScene SceneId = 6
)

var SnakeLevels = []SceneId{
	Level_0,
	Level_1,
	Level_2,
	Level_3,
}

type SceneData struct {
	CurrentScene SceneId
	Project      *assets.LDtkProject
}

var SceneState = donburi.NewComponentType[SceneData]()
