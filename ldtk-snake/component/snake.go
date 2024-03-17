package component

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/engine"
	"github.com/yohamta/donburi"
)

type SnakeBodyData struct {
	Entry    *donburi.Entry
	Next     *SnakeBodyData
	Previous *SnakeBodyData
	History  []HistoryData
}

var SnakeBody = donburi.NewComponentType[SnakeBodyData]()

type HistoryData struct {
	Position resolv.Vector
	Velocity resolv.Vector
}
type SnakeData struct {
	Speed        float64
	History      []HistoryData
	HistoryTimer *engine.Timer
	Tail         *SnakeBodyData
}

var Snake = donburi.NewComponentType[SnakeData]()

func GetTail(snake *SnakeData) (*SnakeBodyData, *SnakeBodyData) {
	next := snake.Tail

	if snake.Tail == nil {
		return nil, snake.Tail
	}
	var prev *SnakeBodyData
	for next != nil {
		prev = next
		next = next.Next
	}
	return prev, next
}

func (sd *SnakeData) SetTail(part *SnakeBodyData) {
	if sd.Tail == nil {
		sd.Tail = part
		return
	}

	next := sd.Tail
	for next.Next != nil {
		next = next.Next
	}
	next.Next = part
}
