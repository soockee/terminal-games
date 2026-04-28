//go:build !autoplay

package system

import "github.com/soockee/terminal-games/match-3/component"

// tryAutoPlay is a no-op when built without the "autoplay" tag.
func tryAutoPlay(_ *component.BoardData) bool {
	return false
}

const autoPlayEnabled = false
