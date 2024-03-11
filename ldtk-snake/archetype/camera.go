package archetype

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

var Camera = newArchetype(
	transform.Transform,
	component.Velocity,
	component.Camera,
)

func MustFindCamera(w donburi.World) *donburi.Entry {
	camera, ok := query.NewQuery(filter.Contains(component.Camera)).First(w)
	if !ok {
		panic("no camera found")
	}

	return camera
}
