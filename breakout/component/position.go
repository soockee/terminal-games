package component

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

var InputPosition = donburi.NewComponentType[input.Vec]()
var SpacePosition = donburi.NewComponentType[resolv.Vector]()
