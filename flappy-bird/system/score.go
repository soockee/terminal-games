package system

import (
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateScore awards a point when the bird passes a pipe pair.
// Each pipe entity is marked Passed once the bird's left edge exceeds
// the pipe's right edge. Since each pair has two pipe entities (top+bottom),
// we award half a point per entity (i.e. one point per pair) by only
// counting pipes where the entry also has PipeData.Passed == false.
func UpdateScore(e *ecs.ECS) {
	birdEntry, ok := tags.Bird.First(e.World)
	if !ok {
		return
	}
	scoreEntry, ok := component.Score.First(e.World)
	if !ok {
		return
	}
	goEntry, ok := component.GameOver.First(e.World)
	if !ok {
		return
	}
	gd := component.GameOver.Get(goEntry)
	if gd.Dead || gd.Won {
		return
	}

	birdX := component.Shape.Get(birdEntry).Shape.Bounds().Min.X
	sd := component.Score.Get(scoreEntry)

	tags.Pipe.Each(e.World, func(entry *donburi.Entry) {
		pd := component.Pipe.Get(entry)
		if pd.Passed {
			return
		}
		pipeRight := component.Shape.Get(entry).Shape.Bounds().Max.X
		if birdX > pipeRight {
			pd.Passed = true
			// Two pipe entities per pair — only count the top (flipped) pipe.
			if pd.FlipY {
				sd.Value++
			}
		}
	})

	// Win condition: all pipe pairs passed.
	if sd.Target > 0 && sd.Value >= sd.Target {
		gd.Won = true
	}
}
