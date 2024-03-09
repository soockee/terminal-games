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

type InputOverlayType int

const (
	UpField InputOverlayType = iota
	DownField
	LeftField
	RightField
)

type Direction CellPosition

var (
	CellTypeMapping = map[CellType]color.Color{
		SnakeHead: color.RGBA{255, 0, 0, 255},
		SnakeBody: color.RGBA{100, 50, 0, 255},
		Food:      color.RGBA{190, 0, 190, 255},
		EmptyCell: color.RGBA{255, 255, 255, 0},
		Wall:      color.RGBA{0, 0, 0, 255},
	}
	OverlayTypeMapping = map[InputOverlayType]color.Color{
		UpField:    color.RGBA{17, 66, 50, 128},
		DownField:  color.RGBA{17, 66, 50, 128},
		LeftField:  color.RGBA{135, 169, 34, 255},
		RightField: color.RGBA{92, 131, 116, 255},
	}
)
