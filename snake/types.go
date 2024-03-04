package main

import "image/color"

type CellPosition struct {
	dx int
	dy int
}

func (cp CellPosition) equals(cmp CellPosition) bool {
	if cp.dx == cmp.dx && cp.dy == cmp.dy {
		return true
	}
	return false
}

type Cell struct {
	pos      CellPosition
	cellType CellType
}

func NewCell(pos CellPosition, cellType CellType) *Cell {
	return &Cell{
		pos:      pos,
		cellType: cellType,
	}
}

type CellType int

const (
	SnakeHead CellType = iota
	SnakeBody
	Food
	EmptyCell
	Wall
)

type Direction CellPosition

var (
	CellTypeMapping = map[CellType]color.Color{
		SnakeHead: color.RGBA{255, 0, 0, 255},
		SnakeBody: color.RGBA{100, 50, 0, 255},
		Food:      color.RGBA{0, 255, 0, 255},
		EmptyCell: color.RGBA{255, 255, 255, 255},
		Wall:      color.RGBA{0, 0, 0, 255},
	}
)
