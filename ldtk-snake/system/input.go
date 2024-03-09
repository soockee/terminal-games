package system

import input "github.com/quasilyte/ebitengine-input"

var InputSytem *input.System

func init() {
	InputSytem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
}
