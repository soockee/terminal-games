package assets

import (
	"image"
	"image/color"
	"image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	imageCache   = map[string]*ebiten.Image{}
	imageCacheMu sync.Mutex
)

// chromaKey is the blue background color used in enemy spritesheets.
var chromaKey = color.RGBA{R: 27, G: 89, B: 153, A: 255}

// LoadImage loads a PNG from the embedded FS and returns a cached *ebiten.Image.
// If replaceChroma is true, the blue background (27,89,153) is replaced with
// full transparency.
func LoadImage(path string, replaceChroma bool) *ebiten.Image {
	key := path
	if replaceChroma {
		key += ":chroma"
	}

	imageCacheMu.Lock()
	defer imageCacheMu.Unlock()

	if cached, ok := imageCache[key]; ok {
		return cached
	}

	f, err := FS.Open(path)
	if err != nil {
		panic("assets: open " + path + ": " + err.Error())
	}
	defer f.Close()

	src, err := png.Decode(f)
	if err != nil {
		panic("assets: decode " + path + ": " + err.Error())
	}

	var result image.Image = src

	if replaceChroma {
		bounds := src.Bounds()
		dst := image.NewNRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := src.At(x, y).RGBA()
				if uint8(r>>8) == chromaKey.R && uint8(g>>8) == chromaKey.G && uint8(b>>8) == chromaKey.B {
					dst.SetNRGBA(x, y, color.NRGBA{})
				} else {
					dst.Set(x, y, src.At(x, y))
				}
			}
		}
		result = dst
	}

	img := ebiten.NewImageFromImage(result)
	imageCache[key] = img
	return img
}
