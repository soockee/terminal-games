package component

import (
	"github.com/yohamta/donburi"
)

type TextData struct {
	Identifier          string
	Text                []string
	IsAnimated          bool
	CurrentCharPosition int
}

var Text = donburi.NewComponentType[TextData]()
