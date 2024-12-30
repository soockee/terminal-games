package assets

import (
	"bytes"
	_ "embed"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed fonts/kenney_pixel.ttf
	NormalFontData []byte
	//go:embed fonts/kenney_pixel_square.ttf
	SquareFontData []byte

	SmallFont  *text.GoTextFace
	NormalFont *text.GoTextFace
	SqaureFont *text.GoTextFace
)

func MustLoadAssets() {
	SmallFont = MustLoadFont(NormalFontData, 10)
	NormalFont = MustLoadFont(NormalFontData, 24)
	SqaureFont = MustLoadFont(SquareFontData, 24)
}

func MustLoadFont(data []byte, size int) *text.GoTextFace {

	s, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	f := &text.GoTextFace{
		Source: s,
		Size:   float64(size),
	}

	return f
}
