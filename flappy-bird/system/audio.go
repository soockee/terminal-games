package system

import (
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/event"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

// SubscribeAudioEvents wires the audio event handler to the ECS world.
// Call this once after Build() creates the world.
func SubscribeAudioEvents(w donburi.World) {
	event.AudioEvent.Subscribe(w, onAudioEvent)
}

// ProcessEvents drains all pending events and dispatches to subscribers.
func ProcessEvents(e *ecs.ECS) {
	events.ProcessAllEvents(e.World)
}

func onAudioEvent(w donburi.World, e event.AudioEventData) {
	entry, ok := component.Audio.First(w)
	if !ok {
		return
	}
	ad := component.Audio.Get(entry)
	player, ok := ad.SFX[e.Name]
	if !ok {
		return
	}
	// Rewind + replay: one active instance per SFX.
	if err := player.Rewind(); err == nil {
		player.Play()
	}
}
