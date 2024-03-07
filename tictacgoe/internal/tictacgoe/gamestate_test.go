package tictacgoe

import "testing"

func TestIsWin(t *testing.T) {
	gameState := NewGamestate()

	for i := 0; i < 5; i++ {
		gameState.moveCounter = i
		if gameState.IsWin() == true {
			t.Error("Game cannot be won before 5th move")
		}
	}

	gameState.moveCounter = 5
	/**
	* -----
	* |xxx|
	* |   |
	* |   |
	* -----
	*
	 */
	gameState.boards[0] = 0b000_000_000_000_000_000_000_111
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}

	/**
	* -----
	* |   |
	* |xxx|
	* |   |
	* -----
	*
	 */
	gameState.boards[0] = 0b000_000_000_000_000_000_111_000
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}

	/**
	* -----
	* |   |
	* |   |
	* |xxx|
	* -----
	*
	 */
	gameState.boards[0] = 0b000_000_000_000_000_111_000_000
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}

	/**
	* -----
	* |x  |
	* |x  |
	* |x  |
	* -----
	*
	 */
	gameState.boards[0] = 0b000_000_000_000_111_000_000_000
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}

	/**
	* -----
	* | x |
	* | x |
	* | x |
	* -----
	*
	 */
	gameState.boards[0] = 0b000_000_000_111_000_000_000_000
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}

	/**
	* -----
	* |  x|
	* |  x|
	* |  x|
	* -----
	*
	 */
	gameState.boards[0] = 0b000_000_111_000_000_000_000_000
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}

	/**
	* -----
	* |x  |
	* | x |
	* |  x|
	* -----
	*
	 */
	gameState.boards[0] = 0b000_111_000_000_000_000_000_000
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}

	/**
	* -----
	* |  X|
	* | x |
	* |X  |
	* -----
	*
	 */
	gameState.boards[0] = 0b111_000_000_000_000_000_000_000
	if gameState.IsWin() != true {
		t.Error("Game should be won")
	}
}

func TestMakeMove(t *testing.T) {
	gameState := NewGamestate()

	gameState.MakeMove(0)
	if gameState.getCurrentPlayer() != 1 {
		t.Error("current palyer did not change")
	}
	if gameState.boards[0] != 0b000_001_000_000_001_000_000_001 {
		t.Error("board was not moved by move 0")
	}

	v, _ := gameState.history.Peek()
	if v != 0 {
		t.Error("wrong move saved")
	}
}

func TestMakeMoves(t *testing.T) {
	gameState := NewGamestate()
	gameState.MakeMoves(0, 1)
	if gameState.boards[0] != 0b000_001_000_000_001_000_000_001 {
		t.Error("board was not moved by move 0")
	}
	if gameState.boards[1] != 0b000_000_000_001_000_000_000_010 {
		t.Error("board was not moved by move 1")
	}
	gameState.MakeMove(2)
	if gameState.boards[0] != 0b001_001_001_000_001_000_000_101 {
		t.Error("board was not moved by move 2")
	}
}

func TestUndoMove(t *testing.T) {
	gameState := NewGamestate()
	gameState.MakeMove(0)
	gameState.UndoMove()
	if gameState.boards[0] != 0b000_000_000_000_000_000_000_000 {
		t.Error("move was not undone")
	}
	_, err := gameState.history.Peek()
	if err == nil {
		t.Error("history is incorrect")
	}
}

func TestBoardToHash(t *testing.T) {
	gameState := NewGamestate()
	gameState.MakeMove(0)
	hash := gameState.BoardToHash()
	if hash != 0b000_000_001_000_000_000 {
		t.Error("hashing failed")
	}

	gameState.MakeMove(1)
	hash2 := gameState.BoardToHash()
	if hash2 != 0b_000_000_001_000_000_010 {
		t.Error("hashing failed")
	}
}

// func TestFlipDiagonal(t *testing.T) {
// }

// func TestFlipAntiDiagonal(t *testing.T) {
// }

// func TestFlipHorizontal(t *testing.T) {
// }

// func TestFlipVertical(t *testing.T) {
// }

// func TestRotate90DegreesClockwise(t *testing.T) {
// }

// func TestRotate90DegreesAntiClockwise(t *testing.T) {
// }

// func TestRotate180Degrees(t *testing.T) {
// }

// func TestRotate90DegreesClockwiseMirror(t *testing.T) {
// }

// func TestRotate90DegreesAntiClockWiseMirror(t *testing.T) {
// }
