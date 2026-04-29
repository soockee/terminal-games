//go:build !autoplay

package system

import "github.com/soockee/terminal-games/match-3/component"

// tryAutoPlay is a no-op when built without the "autoplay" tag.
func tryAutoPlay(_ *component.GridData, _ *component.PhaseData, _ *component.InputData) bool {
	return false
}

const autoPlayEnabled = false
