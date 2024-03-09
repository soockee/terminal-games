package systems

import input "github.com/quasilyte/ebitengine-input"

var InputSytem *input.System = &input.System{}

func init() {
	InputSytem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
}
