package util

import "math/rand"

var RandomNames = []string{
	"Otto Obst",
	"Günther",
	"Ilse",
	"Klaus",
	"Liesel",
	"Uwe",
	"Wanda",
	"Franz",
	"Max",
	"Karl",
	"Emil",
	"Hans",
	"Dieter",
	"Helga",
	"Inge",
	"Kurt",
	"Lisa",
	"Moritz",
	"Nora",
	"Paula",
	"Rudi",
	"Sepp",
	"Greta",
	"Bruno",
	"Martha",
	"Willi",
}

func GetRandomName() string {
	// Use math/rand to generate a random index within the slice
	randomIndex := rand.Intn(len(RandomNames))
	return RandomNames[randomIndex]
}
