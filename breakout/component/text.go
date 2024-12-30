package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type TextData struct {
	Text       []string
	IsAnimated bool
	Shape      resolv.IShape
}

var Text = donburi.NewComponentType[TextData]()
