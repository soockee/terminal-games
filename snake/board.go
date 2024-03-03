package main

import (
	"image/color"
)

const (
	CellsDX  = 64
	CellsDY  = 64
	GridSize = 16
	Offset   = 2
)

type Board struct {
	cells [][]Cell
	snake *Snake
	ticks int
}

type Cell int

const (
	SnakeHead Cell = iota
	SnakeBody
	Food
	EmptyCell
	Wall
)

var (
	CellMapping = map[Cell]color.Color{
		SnakeHead: color.RGBA{255, 0, 0, 255},
		SnakeBody: color.RGBA{153, 0, 0, 255},
		Food:      color.RGBA{51, 204, 51, 255},
		EmptyCell: color.RGBA{255, 255, 255, 255},
		Wall:      color.RGBA{0, 0, 0, 255},
	}
)

func NewBoard() *Board {
	board := Board{
		cells: make([][]Cell, CellsDX),
		snake: NewSnake(),
	}
	for i := range board.cells {
		board.cells[i] = make([]Cell, CellsDY)
	}

	for dx, row := range board.cells {
		for dy := range row {
			if dy == 0 || dy == CellsDY-1 {
				board.cells[dx][dy] = Wall
			} else if dx == 0 || dx == CellsDX-1 {
				board.cells[dx][dy] = Wall
			} else {
				board.cells[dx][dy] = EmptyCell
			}
		}
	}

	board.cells[board.snake.snakeHeadDX][board.snake.snakeHeadDY] = SnakeHead

	return &board
}

func (b *Board) UpdateActors() {
	b.snake.CheckDirection()
	if b.ticks > b.snake.snakeSpeed {
		b.snake.direction(b.snake)
		b.ticks = 0
	}
}
