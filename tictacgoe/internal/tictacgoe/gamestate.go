package tictacgoe

/**
* Board:
* 012
* 345
* 678
* 246_048_258_147_036_876_543_210
* boards[0] = 'x'
* boards[1] = 'o'
 */

const (
	bitfilter         = 0b001_001_001_001_001_001_001_001
	bitfilterForBoard = 0b111_111_111
)

var (
	BitPattern = [...]int{
		0b000_001_000_000_001_000_000_001,
		0b000_000_000_001_000_000_000_010,
		0b001_000_001_000_000_000_000_100,
		0b000_000_000_000_010_000_001_000,
		0b010_010_000_010_000_000_010_000,
		0b000_000_010_000_000_000_100_000,
		0b100_000_000_000_100_001_000_000,
		0b000_000_000_100_000_010_000_000,
		0b000_100_100_000_000_100_000_000,
	}
)

type Gamestate struct {
	boards      [2]int
	moveCounter int
	history     *Stack[int]
}

func NewGamestate() *Gamestate {
	return &Gamestate{
		boards: [...]int{
			0b000_000_000_000_000_000_000_000,
			0b000_000_000_000_000_000_000_000,
		},
		moveCounter: 0,
		history:     NewStack[int](),
	}
}

func (gamestate *Gamestate) IsWin() bool {
	if gamestate.moveCounter < 5 {
		return false
	}
	turn := (gamestate.moveCounter - 1) & 1
	currentBoard := gamestate.boards[turn]
	return (currentBoard & (currentBoard >> 1) & (currentBoard >> 2) & bitfilter) > 0
}

func (gamestate *Gamestate) MakeMove(move int) {
	turn := gamestate.moveCounter & 1
	gamestate.history.Push(turn)
	gamestate.boards[turn] |= BitPattern[move]
	gamestate.moveCounter++
}

func (gamestate *Gamestate) MakeMoves(moves ...int) {
	for move := range moves {
		gamestate.MakeMove(move)
	}
}

func (gamestate *Gamestate) UndoMove() {
	gamestate.moveCounter--
	lastMove, err := gamestate.history.Pop()
	if err == nil {
		gamestate.boards[gamestate.moveCounter&1] = lastMove
	}
}

func (gamestate *Gamestate) BoardToHash() int {
	board0ToHash := (gamestate.boards[0] & bitfilterForBoard) << 9
	board1ToHash := gamestate.boards[1] & bitfilterForBoard
	return board0ToHash | board1ToHash
}

func FlipDiagonal(board int) int {
	res := board & bitfilterForBoard
	res = (res << 8) & 100_000_000
	res = res | (res>>8)&0b1
	res = res | (res<<4)&0b100_000
	res = res | (res>>4)&0b10
	res = res | (res<<4)&0b10_000_000
	res = res | (res>>4)&0b1_000
	res = res | (board & 0b001_010_100)
	return res
}

func FlipAntiDiagonal(board int) int {
	res := board & bitfilterForBoard
	res = (res << 4) & 0b1_000_000
	res = res | (res>>4)&0b100
	res = res | (res<<2)&0b10_000_000
	res = res | (res>>2)&0b100_000
	res = res | (res<<2)&0b1_000
	res = res | (res>>2)&0b10
	res = res | (board & 0b100_010_001)
	return res
}

func FlipHorizontal(board int) int {
	res := board & bitfilterForBoard
	res = ((res << 6) & 0b111_000_000) | ((res >> 6) & 0b000_000_111)
	res = res | (board & 0b000_111_000)
	return res
}

func FlipVertical(board int) int {
	res := board & bitfilterForBoard
	res = ((res << 2) & 0b100_100_100) | ((res >> 2) & 0b001_001_001)
	res = res | (board & 0b010_010_010)
	return res
}

func Rotate90DegreesClockwise(board int) int {
	return FlipHorizontal(FlipDiagonal(board))
}

func Rotate90DegreesAntiClockwise(board int) int {
	return FlipHorizontal(FlipAntiDiagonal(board))
}

func Rotate180Degrees(board int) int {
	return FlipHorizontal(FlipVertical(board))
}

func Rotate90DegreesClockwiseMirror(board int) int {
	return Rotate90DegreesClockwise(FlipVertical(board))
}

func Rotate90DegreesAntiClockWiseMirror(board int) int {
	return Rotate90DegreesAntiClockwise(FlipVertical(board))
}

// x starts
// 0 = 'x' ; 1 = 'o';
func (gamestate *Gamestate) getCurrentPlayer() int {
	return gamestate.moveCounter & 1
}
