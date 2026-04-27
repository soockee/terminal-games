package archetype

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/yohamta/donburi"
)

// NewSpace creates a singleton entity holding the resolv.Space.
func NewSpace(w donburi.World, width, height, cellW, cellH int) {
	e := w.Entry(w.Create(component.Space))
	space := resolv.NewSpace(width, height, cellW, cellH)
	component.Space.Set(e, space)
}

// NewDebug creates the singleton debug toggle entity.
func NewDebug(w donburi.World) {
	e := w.Entry(w.Create(component.Debug))
	component.Debug.Set(e, &component.DebugData{Enabled: false})
}
