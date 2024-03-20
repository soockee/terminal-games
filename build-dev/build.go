package main

import (
	"os"
	"os/exec"
)

func main() {
	buildLDtkSnake()
}

func buildLDtkSnake() {
	cmd := exec.Command("go", "build", "-o", "ldtk-snake.wasm", "../ldtk-snake")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Run()
}
