package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

// AnimationData holds a set of named animations and tracks which one is active.
// Use string keys for animation states (e.g. "idle", "walk", "jump", "fall").
type AnimationData struct {
	Animations map[string]*ganim8.Animation
	Current    string
	FlipH      bool // mirror sprite horizontally (e.g. walking left)
}

var Animation = donburi.NewComponentType[AnimationData]()

// Play switches to the named animation. If the same animation is already
// playing, this is a no-op. When switching, the new animation resets to
// frame 1 and resumes playback.
func (d *AnimationData) Play(name string) {
	if d.Current == name {
		return
	}
	if _, ok := d.Animations[name]; !ok {
		return
	}
	d.Current = name
	d.Animations[name].GoToFrame(1)
	d.Animations[name].Resume()
}

// CurrentAnimation returns the active *ganim8.Animation, or nil.
func (d *AnimationData) CurrentAnimation() *ganim8.Animation {
	if d.Current == "" {
		return nil
	}
	return d.Animations[d.Current]
}
