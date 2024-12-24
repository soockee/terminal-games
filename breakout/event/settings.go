package event

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi/features/events"
)

type UpdateSetting struct {
	Action input.Action
}

var UpdateSettingEvent = events.NewEventType[*UpdateSetting]()
