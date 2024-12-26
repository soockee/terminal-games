package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type TextData struct {
	Identifier          string
	Text                []string
	IsAnimated          bool
	CurrentCharPosition int
	Shape               *resolv.ConvexPolygon
}

var Text = donburi.NewComponentType[TextData]()
