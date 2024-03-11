package component

import (
	"github.com/yohamta/donburi"
)

type ButtonData struct {
	Clicked     bool
	HandlerFunc func()
}

var Button = donburi.NewComponentType[ButtonData]()
