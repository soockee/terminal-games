package component

import (
	"github.com/yohamta/donburi"
)

type ButtonData struct {
	Clicked     bool
	HandlerFunc func(donburi.World)
}

var Button = donburi.NewComponentType[ButtonData]()
