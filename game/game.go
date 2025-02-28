package game

import (
	"fmt"
)

type u3t struct {
	XBoard      uint64
	OBoard      uint64
	ActiveBoard int
	XTurn       bool
}

const (
	EmptyBoard uint64 = 0
	AllCells   uint64 = 0x1FFFFFFFF
)

func NewGame() *u3t {
	return &u3t{
		XBoard:      EmptyBoard,
		OBoard:      EmptyBoard,
		ActiveBoard: -1,
		XTurn:       true, // First move is played by X
	}
}

// GetBitPosition converts row, col in global board to bit position
// globalRow, globalCol: 0-8 coordinates on the entire 9×9 board
func GetBitPosition(globalRow, globalCol int) int {
	return globalRow*9 + globalCol
}

func GetGlobalPosition(board, localPos int) int {
	boardRow := board / 3
	boardCol := board % 3
	localRow := localPos / 3
	localCol := localPos % 3

	globalRow := boardRow*3 + localRow
	globalCol := boardCol*3 + localCol

	return GetBitPosition(globalRow, globalCol)
}

func (g *u3t) IsEmpty(bitPos int) bool {
	mask := uint64(1) << bitPos
	return ((g.XBoard & mask) == 0) && ((g.OBoard & mask) == 0)
}

func (g *u3t) MakeMove(board, localPos int) bool {
	if g.ActiveBoard != -1 && g.ActiveBoard != board {
		return false
	}

	bitPos := GetGlobalPosition(board, localPos)
	mask := uint64(1) << bitPos

	if !g.IsEmpty(bitPos) {
		return false
	}

	if g.XTurn {
		g.XBoard |= mask
	} else {
		g.OBoard |= mask
	}

	g.ActiveBoard = localPos

	if g.IsBoardFull(g.ActiveBoard) || g.IsBoardWon(g.ActiveBoard) {
		g.ActiveBoard = -1
	}

	// Switch turns
	g.XTurn = !g.XTurn

	return true
}

// IsBoardWon checks if a small board is won
func (g *u3t) IsBoardWon(board int) bool {
	// Convert small board to a 3×3 representation for easier win checking
	var xSmallBoard, oSmallBoard uint16

	for i := range 9 {
		bitPos := GetGlobalPosition(board, i)
		mask := uint64(1) << bitPos

		if (g.XBoard & mask) != 0 {
			xSmallBoard |= (1 << i)
		} else if (g.OBoard & mask) != 0 {
			oSmallBoard |= (1 << i)
		}
	}

	// Win patterns for a 3×3 board
	winPatterns := []uint16{
		0x7,   // Row 1: 111 000 000
		0x38,  // Row 2: 000 111 000
		0x1C0, // Row 3: 000 000 111
		0x49,  // Col 1: 001 001 001
		0x92,  // Col 2: 010 010 010
		0x124, // Col 3: 100 100 100
		0x111, // Diag 1: 100 010 001
		0x54,  // Diag 2: 001 010 100
	}

	// Check if any win pattern is satisfied
	for _, pattern := range winPatterns {
		if (xSmallBoard & pattern) == pattern {
			return true // X won this small board
		}
		if (oSmallBoard & pattern) == pattern {
			return true // O won this small board
		}
	}

	return false
}

// IsBoardFull checks if a small board is full
func (g *u3t) IsBoardFull(board int) bool {
	for i := range 9 {
		bitPos := GetGlobalPosition(board, i)
		if g.IsEmpty(bitPos) {
			return false // Found an empty cell, board is not full
		}
	}
	return true
}

// GetBoardWinner returns 1 for X, 2 for O, 0 for no winner
func (g *u3t) GetBoardWinner(board int) int {
	if !g.IsBoardWon(board) {
		return 0
	}

	var xSmallBoard, oSmallBoard uint16

	for i := range 9 {
		bitPos := GetGlobalPosition(board, i)
		mask := uint64(1) << bitPos

		if (g.XBoard & mask) != 0 {
			xSmallBoard |= (1 << i)
		} else if (g.OBoard & mask) != 0 {
			oSmallBoard |= (1 << i)
		}
	}

	winPatterns := []uint16{
		0x7, 0x38, 0x1C0, 0x49, 0x92, 0x124, 0x111, 0x54,
	}

	for _, pattern := range winPatterns {
		if (xSmallBoard & pattern) == pattern {
			return 1 // X won
		}
		if (oSmallBoard & pattern) == pattern {
			return 2 // O won
		}
	}

	return 0 // Should never reach here if IsBoardWon is true
}

// IsGameWon checks if the whole game is won
func (g *u3t) IsGameWon() bool {
	// Create a 3×3 representation of the meta-board
	var xMetaBoard, oMetaBoard uint16

	for board := range 9 {
		winner := g.GetBoardWinner(board)
		if winner == 1 {
			xMetaBoard |= (1 << board)
		} else if winner == 2 {
			oMetaBoard |= (1 << board)
		}
	}

	// Use the same win patterns to check the meta-board
	winPatterns := []uint16{
		0x7,   // Row 1
		0x38,  // Row 2
		0x1C0, // Row 3
		0x49,  // Col 1
		0x92,  // Col 2
		0x124, // Col 3
		0x111, // Diag 1
		0x54,  // Diag 2
	}

	for _, pattern := range winPatterns {
		if (xMetaBoard & pattern) == pattern {
			return true // X won the game
		}
		if (oMetaBoard & pattern) == pattern {
			return true // O won the game
		}
	}

	return false
}

// GetGameWinner returns 1 for X, 2 for O, 0 for no winner
func (g *u3t) GetGameWinner() int {
	// Create a 3×3 representation of the meta-board
	var xMetaBoard, oMetaBoard uint16

	for board := range 9 {
		winner := g.GetBoardWinner(board)
		if winner == 1 {
			xMetaBoard |= (1 << board)
		} else if winner == 2 {
			oMetaBoard |= (1 << board)
		}
	}

	// Use the same win patterns to check the meta-board
	winPatterns := []uint16{
		0x7, 0x38, 0x1C0, 0x49, 0x92, 0x124, 0x111, 0x54,
	}

	for _, pattern := range winPatterns {
		if (xMetaBoard & pattern) == pattern {
			return 1 // X won the game
		}
		if (oMetaBoard & pattern) == pattern {
			return 2 // O won the game
		}
	}

	return 0 // No winner yet
}

// IsGameDraw checks if the game is a draw
func (g *u3t) IsGameDraw() bool {
	// Game is a draw if all cells are filled and no one has won
	for i := range 81 {
		if g.IsEmpty(i) {
			return false // Found an empty cell, game is not a draw
		}
	}
	return !g.IsGameWon()
}

// PrintBoard prints the current state of the board
func (g *u3t) PrintBoard() {
	board := make([][]rune, 9)
	for i := range board {
		board[i] = make([]rune, 9)
	}

	// Fill the board representation
	for row := range 9 {
		for col := range 9 {
			bitPos := GetBitPosition(row, col)
			mask := uint64(1) << bitPos

			if (g.XBoard & mask) != 0 {
				board[row][col] = 'X'
			} else if (g.OBoard & mask) != 0 {
				board[row][col] = 'O'
			} else {
				board[row][col] = '.'
			}
		}
	}

	// Print the board with separators
	fmt.Println("Ultimate Tic-Tac-Toe")
	fmt.Println("===================")

	for row := 0; row < 9; row++ {
		if row > 0 && row%3 == 0 {
			fmt.Println("---+---+---")
		}

		for col := 0; col < 9; col++ {
			if col > 0 && col%3 == 0 {
				fmt.Print("|")
			}
			fmt.Printf("%c", board[row][col])
		}
		fmt.Println()
	}

	// Print active board information
	if g.ActiveBoard == -1 {
		fmt.Println("Next player can choose any board")
	} else {
		fmt.Printf("Next move must be in board %d\n", g.ActiveBoard)
	}

	if g.XTurn {
		fmt.Println("X's turn")
	} else {
		fmt.Println("O's turn")
	}
}
