package tictacgoe

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type PerformanceInfo struct{}

func NewPerformanceInfo() *PerformanceInfo {
	performanceInfo := &PerformanceInfo{}
	return performanceInfo
}

const (
	PerformanceInfoWidth  = 115
	PerformanceInfoHeight = 50
	PerformanceInfoMargin = 25
)

func (performanceInfo *PerformanceInfo) Update() error {
	return nil
}

func (performanceInfo *PerformanceInfo) Draw(performanceInfoImage *ebiten.Image, font font.Face) {
	const x = 0
	performanceInfoImage.Clear()
	// Draw info
	msg := fmt.Sprintf("TPS %0.2f", ebiten.ActualTPS())
	text.Draw(performanceInfoImage, msg, font, x, 24, color.Black)
}

func (performanceInfo *PerformanceInfo) Size() (int, int) {
	return PerformanceInfoWidth, PerformanceInfoHeight
}
func (performanceInfo *PerformanceInfo) Margin() int {
	return PerformanceInfoMargin
}
