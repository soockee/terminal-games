package component

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi"
)

var Position = donburi.NewComponentType[input.Vec]()
