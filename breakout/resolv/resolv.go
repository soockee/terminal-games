package resolv

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
)

func Add(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		component.Space.Get(space).Add(GetObject(obj))
	}
}

func Remove(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		component.Space.Get(space).Remove(GetObject(obj))
	}
}

func SetObject(entry *donburi.Entry, obj *resolv.ConvexPolygon) {
	component.ConvexPolygon.Set(entry, obj)
}

func GetObject(entry *donburi.Entry) *resolv.ConvexPolygon {
	return component.ConvexPolygon.Get(entry)
}
