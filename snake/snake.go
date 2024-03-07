package main

import (
	"log/slog"

	input "github.com/quasilyte/ebitengine-input"
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

	direction  input.Action
	snakeSpeed int
	input      *input.Handler
}

func NewSnake(input *input.Handler) *Snake {
	snake := &Snake{
		snakeHead:  newSnakeHead(randomPosition()),
		snakeSpeed: 2,
		input:      input,
	}
	return snake
}

func (s *Snake) CheckDirection() {
	if s.input.ActionIsPressed(ActionMoveUp) {
		slog.Debug("Up")
		s.direction = ActionMoveUp
	}
	if s.input.ActionIsPressed(ActionMoveDown) {
		slog.Debug("Down")
		s.direction = ActionMoveDown
	}
	if s.input.ActionIsPressed(ActionMoveLeft) {
		slog.Debug("Left")
		s.direction = ActionMoveLeft
	}
	if s.input.ActionIsPressed(ActionMoveRight) {
		slog.Debug("Right")
		s.direction = ActionMoveRight
	}
	if info, ok := s.input.JustPressedActionInfo(ActionClick); ok {
		if info.Pos.Y < fieldHeight {
			slog.Debug("Up")
			s.direction = ActionMoveUp
		} else if info.Pos.Y < fieldHeight*3 && info.Pos.Y > fieldHeight*2 {
			slog.Debug("Down")
			s.direction = ActionMoveDown
		} else if info.Pos.Y < fieldHeight*2 && info.Pos.Y > fieldHeight && info.Pos.X < fieldWidth {
			slog.Debug("Left")
			s.direction = ActionMoveLeft
		} else if info.Pos.Y < fieldHeight*2 && info.Pos.Y > fieldHeight && info.Pos.X > fieldWidth {
			slog.Debug("Right")
			s.direction = ActionMoveRight
		}
	}
}

func (s *Snake) move(action input.Action) (head *CellPosition, tail *CellPosition) {
	head = &CellPosition{s.snakeHead.pos.dx, s.snakeHead.pos.dy}
	switch action {
	case ActionMoveUp:
		s.snakeHead.pos.dy--
	case ActionMoveDown:
		s.snakeHead.pos.dy++
	case ActionMoveLeft:
		s.snakeHead.pos.dx--
	case ActionMoveRight:
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
