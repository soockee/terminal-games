package system

import (
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	collectionPlayerQuery = donburi.NewQuery(
		filter.Contains(component.Player, component.Body),
	)
	collectableQuery = donburi.NewQuery(
		filter.Contains(component.Collectable, component.Position, component.Sprite),
	)
)

// UpdateCollection checks for overlap between the player and any collectable
// entity each tick. Collected entities are removed and score is incremented.
func UpdateCollection(e *ecs.ECS) {
	gsEntry, ok := component.GameState.First(e.World)
	if !ok || !component.GameState.Get(gsEntry).IsActive() {
		return
	}

	var px, py, pw, ph float64
	collectionPlayerQuery.Each(e.World, func(entry *donburi.Entry) {
		body := component.Body.Get(entry).Body
		if body == nil || !body.IsAlive() {
			return
		}
		px, py = body.Position()
		pw, ph = body.Size()
	})
	if pw == 0 {
		return
	}

	scoreEntry, hasScore := component.Score.First(e.World)

	var toCollect []*donburi.Entry
	collectableQuery.Each(e.World, func(entry *donburi.Entry) {
		pos := component.Position.Get(entry)
		spr := component.Sprite.Get(entry)
		if spr.Image == nil {
			return
		}
		b := spr.Image.Bounds()
		if overlapsAABB(px, py, pw, ph, pos.X, pos.Y, float64(b.Dx()), float64(b.Dy())) {
			toCollect = append(toCollect, entry)
		}
	})

	for _, entry := range toCollect {
		if hasScore {
			component.Score.Get(scoreEntry).Value += 100
		}
		entry.Remove()
	}
}

func overlapsAABB(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}
