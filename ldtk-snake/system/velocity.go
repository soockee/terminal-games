package system

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type Velocity struct {
	query *query.Query
}

func NewVelocity() *Velocity {
	return &Velocity{
		query: query.NewQuery(
			filter.Contains(transform.Transform, component.Velocity),
		),
	}
}

func (v *Velocity) Update(w donburi.World) {
	v.query.Each(w, func(entry *donburi.Entry) {
		snakeObject := dresolv.GetObject(entry)
		velocity := component.Velocity.Get(entry)
		snakeObject.Position = snakeObject.Position.Add(velocity.Velocity)
	})
}
