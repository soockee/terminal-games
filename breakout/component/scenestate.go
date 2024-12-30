package component

import (
	"math/rand"
	"reflect"
	"slices"

	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/yohamta/donburi"
	"golang.org/x/exp/maps"
)

const (
	Empty   = "Level_Empty"
	Level_0 = "Level_0"
	Level_1 = "Level_1"

	StartScene      = "Level_Start"
	LevelClearScene = "Level_Clear"
	GameOverScene   = "Level_GameOver"
)

// key naming matters, determines order of levels
var Levels = map[string]bool{
	Level_0: false,
	Level_1: false,
}

type SceneData struct {
	CurrentScene string
	LastScene    string
	Project      *assets.LDtkProject
	Levels       map[string]bool
}

var SceneState = donburi.NewComponentType[SceneData]()

func GetRandomUnplayedLevel() (string, bool) {
	var randomLevel string
	foundUnplayed := false

	for !foundUnplayed {
		randomIndex := rand.Intn(len(Levels))
		randomLevel = string(reflect.ValueOf(Levels).MapKeys()[randomIndex].Interface().(string))
		foundUnplayed = !Levels[randomLevel]
	}

	return randomLevel, foundUnplayed
}

func GetNextLevel(current string) (string, bool) {
	levels := maps.Keys(Levels)
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
