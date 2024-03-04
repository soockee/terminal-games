package main

import (
	"math/rand"
)

const (
	CellsDX  = 64
	CellsDY  = 64
	GridSize = 16
	Offset   = 2
)

type Board struct {
	cells   [][]*Cell
	snake   *Snake
	food    CellPosition
	ticks   int
	hitWall bool
}

func NewBoard() *Board {
	board := Board{
		cells:   make([][]*Cell, CellsDX),
		snake:   NewSnake(),
		food:    randomPosition(),
		hitWall: false,
	}
	for i := range board.cells {
		board.cells[i] = make([]*Cell, CellsDY)
	}

	for dx, row := range board.cells {
		for dy := range row {
			if dy == 0 || dy == CellsDY-1 {
				board.cells[dx][dy] = NewCell(CellPosition{dx: dx, dy: dy}, Wall)
			} else if dx == 0 || dx == CellsDX-1 {
				board.cells[dx][dy] = NewCell(CellPosition{dx: dx, dy: dy}, Wall)
			} else {
				board.cells[dx][dy] = NewCell(CellPosition{dx: dx, dy: dy}, EmptyCell)
			}
		}
	}
	board.cells[board.food.dx][board.food.dy] = NewCell(board.food, Food)

	return &board
}

func (b *Board) UpdateActors() {
	b.snake.CheckDirection()
	if b.ticks > b.snake.snakeSpeed {
		head, tail := b.snake.move(b.snake.direction)
		if head.equals(b.food) {
			b.snake.appendBody(*tail)
			b.cells[head.dx][head.dy] = NewCell(*head, SnakeBody)
			b.food = randomPosition()
			for b.cells[b.food.dx][b.food.dy].cellType != EmptyCell {
				b.food = randomPosition()
			}
			b.cells[b.food.dx][b.food.dy] = NewCell(b.food, Food)
		}
		if b.cells[b.snake.snakeHead.pos.dx][b.snake.snakeHead.pos.dy].cellType == Wall {
			b.hitWall = true
			return
		}
		if b.snake.snakeHead.next == nil {
			b.cells[head.dx][head.dy] = NewCell(*head, EmptyCell)
		}
		if tail != nil && b.snake.snakeHead.next != nil {
			b.cells[tail.dx][tail.dy] = NewCell(*tail, EmptyCell)
		}

		b.cells[b.snake.snakeHead.pos.dx][b.snake.snakeHead.pos.dy] = NewCell(b.snake.snakeHead.pos, SnakeHead)
		b.updateBody(b.snake.snakeHead.next)

		b.ticks = 0
	}
}

func (b *Board) updateBody(next *snakeBody) {
	if next == nil {
		return
	}
	b.cells[next.pos.dx][next.pos.dy] = NewCell(next.pos, SnakeBody)
	b.updateBody(next.next)
}

func randomPosition() CellPosition {
	randomDX := rand.Intn((CellsDX - 2)) + 1
	randomDY := rand.Intn((CellsDY - 2)) + 1

	return CellPosition{
		dx: randomDX,
		dy: randomDY,
	}
}
