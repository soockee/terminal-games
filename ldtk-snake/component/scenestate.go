package component

import (
	"math/rand"
	"reflect"
	"slices"

	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/yohamta/donburi"
	"golang.org/x/exp/maps"
)

const (
	Empty   = "Level_Empty"
	Level_0 = "Level_0"
	Level_1 = "Level_1"
	Level_2 = "Level_2"
	Level_3 = "Level_3"
	Level_4 = "Level_4"

	StartScene      = "Level_Start"
	LevelClearScene = "Level_Clear"
	GameOverScene   = "Level_GameOver"
)

// key naming matters, determines order of levels
var SnakeLevels = map[string]bool{
	Level_0: false,
	Level_1: false,
	Level_2: false,
	Level_3: false,
	Level_4: false,
}

type SceneData struct {
	CurrentScene string
	LastScene    string
	Project      *assets.LDtkProject
	SnakeLevels  map[string]bool
}

var SceneState = donburi.NewComponentType[SceneData]()

func GetRandomUnplayedLevel() (string, bool) {
	var randomLevel string
	foundUnplayed := false

	for !foundUnplayed {
		randomIndex := rand.Intn(len(SnakeLevels))
		randomLevel = string(reflect.ValueOf(SnakeLevels).MapKeys()[randomIndex].Interface().(string))
		foundUnplayed = !SnakeLevels[randomLevel]
	}

	return randomLevel, foundUnplayed
}

func GetNextLevel(current string) (string, bool) {
	levels := maps.Keys(SnakeLevels)
	slices.Sort(levels)

	i := slices.Index(levels, current)

	if i < 0 {
		return "", false
	}
	if i+1 >= len(levels) {
		return "", false
	}

	return levels[i+1], true
}
