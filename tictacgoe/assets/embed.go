package assets

import (
	_ "embed"
)

var (
	//go:embed onion_boy_300x300.png
	OnionBoy_png []byte

	//go:embed pig_300x300.png
	Pig_png []byte

	//go:embed handshake_445_295.gif
	Handshake_gif []byte
)
