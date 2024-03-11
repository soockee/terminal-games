package component

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi"
)

var PositionComponent = donburi.NewComponentType[input.Vec]()
