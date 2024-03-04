package main

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type snakeHead struct {
	pos  CellPosition
	next *snakeBody
}
type snakeBody struct {
	pos  CellPosition
	next *snakeBody
}

func newSnakeHead(cellposition CellPosition) *snakeHead {
	return &snakeHead{
		pos: cellposition,
	}
}
func newSnakeBody(cellposition CellPosition) *snakeBody {
	return &snakeBody{
		pos: cellposition,
	}
}

func (s *Snake) appendBody(pos CellPosition) {
	if s.snakeHead.next == nil {
		s.snakeHead.next = newSnakeBody(pos)
		slog.Debug("Append Body to Head")
		return
	}

	current := s.snakeHead.next
	for current.next != nil {
		current = current.next
	}
	current.next = newSnakeBody(pos)
	slog.Debug("Append Body to Body")
}

type Snake struct {
	snakeHead *snakeHead

	direction  ebiten.Key
	snakeSpeed int
}

func NewSnake() *Snake {
	snake := &Snake{}
	snake.snakeHead = newSnakeHead(randomPosition())
	snake.snakeSpeed = 2
	return snake
}

type directions map[ebiten.Key]CellPosition

var possibleMoves = [4]ebiten.Key{ebiten.KeyUp, ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight}

func (s *Snake) CheckDirection() {
	for _, key := range possibleMoves {
		if inpututil.IsKeyJustPressed(key) {
			s.direction = key
		}
	}
}

func (s *Snake) move(k ebiten.Key) (head *CellPosition, tail *CellPosition) {
	head = &CellPosition{s.snakeHead.pos.dx, s.snakeHead.pos.dy}
	switch k {
	case ebiten.KeyUp:
		s.snakeHead.pos.dy--
	case ebiten.KeyDown:
		s.snakeHead.pos.dy++
	case ebiten.KeyLeft:
		s.snakeHead.pos.dx--
	case ebiten.KeyRight:
		s.snakeHead.pos.dx++
	default:
		slog.Warn("invalid move key")
	}

	if s.snakeHead.next != nil {
		delta := s.snakeHead.next.pos.calculcateDelta(*head)
		tail = s.snakeHead.next.updateBody(delta)
	} else {
		tail = head
	}
	return head, tail
}

func (cp1 CellPosition) calculcateDelta(cp2 CellPosition) CellPosition {
	return CellPosition{cp1.dx - cp2.dx, cp1.dy - cp2.dy}
}

func (s *snakeBody) updateBody(delta CellPosition) *CellPosition {
	tail := &CellPosition{s.pos.dx, s.pos.dy}
	s.pos.dx -= 1 * delta.dx
	s.pos.dy -= 1 * delta.dy
	if s.next != nil {
		tail = s.next.updateBody(s.next.pos.calculcateDelta(*tail))
	}
	return tail
}
