package resolv

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
)

func Add(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		component.Space.Get(space).Add(GetObject(obj))
	}
}

func SetObject(entry *donburi.Entry, obj *resolv.Object) {
	component.Object.Set(entry, obj)
}

func GetObject(entry *donburi.Entry) *resolv.Object {
	return component.Object.Get(entry)
}
