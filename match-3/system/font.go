package system

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/soockee/terminal-games/match-3/assets"
)

var fontSource *text.GoTextFaceSource

func init() {
	f, err := assets.FS.Open("ldtk/fonts/font.ttf")
	if err != nil {
		log.Fatalf("open font: %v", err)
	}
	defer f.Close()

	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatalf("parse font: %v", err)
	}
	fontSource = src
}

// FontFace returns a GoTextFace at the given size.
func FontFace(size float64) *text.GoTextFace {
	return &text.GoTextFace{
		Source: fontSource,
		Size:   size,
	}
}
