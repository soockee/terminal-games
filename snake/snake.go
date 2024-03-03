package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/exp/maps"
)

type Snake struct {
	snakeHeadDX int
	snakeHeadDY int

	direction  InputAction
	snakeSpeed int
}

func NewSnake() *Snake {
	snake := &Snake{}
	snake.RandomizeSnakePosition()
	snake.snakeSpeed = 30
	snake.direction = moveRight
	return snake
}

func (s *Snake) RandomizeSnakePosition() {
	randomDX := rand.Intn((CellsDX - 1)) + 1
	randomDY := rand.Intn((CellsDY - 1)) + 1

	s.snakeHeadDX = randomDX
	s.snakeHeadDY = randomDY
}

type inputActions map[ebiten.Key]InputAction

var inputActionMoveMapping = inputActions{
	ebiten.KeyUp:    moveUp,
	ebiten.KeyDown:  moveDown,
	ebiten.KeyLeft:  moveLeft,
	ebiten.KeyRight: moveRight,
}

type InputAction func(*Snake)

func (s *Snake) CheckDirection() {
	for _, key := range maps.Keys(inputActionMoveMapping) {
		if inpututil.IsKeyJustPressed(key) {
			s.direction = inputActionMoveMapping[key]
		}
	}
}

func moveUp(s *Snake) {
	s.snakeHeadDY--
}

func moveDown(s *Snake) {
	s.snakeHeadDY++
}

func moveLeft(s *Snake) {
	s.snakeHeadDX--
}

func moveRight(s *Snake) {
	s.snakeHeadDX++
}
