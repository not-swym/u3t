package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 800
	cellSize     = 60
	boardPadding = 20
)

// Colors
var (
	boardBgColor        = rl.NewColor(242, 242, 242, 255) // Light gray background
	gridLineColor       = rl.NewColor(200, 200, 200, 255) // Light gray grid lines
	mainGridLineColor   = rl.NewColor(100, 100, 100, 255) // Darker lines for main grid
	xColor              = rl.NewColor(231, 76, 60, 255)   // Red for X
	oColor              = rl.NewColor(52, 152, 219, 255)  // Blue for O
	activeBoardColor    = rl.NewColor(241, 196, 15, 50)   // Yellow highlight for active board
	wonBoardColorX      = rl.NewColor(231, 76, 60, 100)   // Semi-transparent red
	wonBoardColorO      = rl.NewColor(52, 152, 219, 100)  // Semi-transparent blue
	textColor           = rl.NewColor(44, 62, 80, 255)    // Dark text color
	statusBgColor       = rl.NewColor(236, 240, 241, 255) // Light background for status
	hoverHighlightColor = rl.NewColor(52, 152, 219, 30)   // Highlight color for hover
)

// Point represents a 2D point
type Point struct {
	X, Y int
}

// UIState holds UI-specific state
type UIState struct {
	Game            *u3t
	WindowWidth     int
	WindowHeight    int
	BoardOffset     Point
	HoveredCell     Point
	HoveredBoard    int
	HoveredPosition int
	GameOver        bool
	Winner          int
	LastMove        Point
}

// NewUIState creates a new UI state
func NewUIState(game *u3t) *UIState {
	return &UIState{
		Game:            game,
		WindowWidth:     screenWidth,
		WindowHeight:    screenHeight,
		BoardOffset:     Point{X: (screenWidth - 9*cellSize) / 2, Y: (screenHeight - 9*cellSize) / 2},
		HoveredCell:     Point{X: -1, Y: -1},
		HoveredBoard:    -1,
		HoveredPosition: -1,
		GameOver:        false,
		Winner:          0,
		LastMove:        Point{X: -1, Y: -1},
	}
}

// RunGame starts the game UI
func RunGame(ui *UIState) {
	// Initialize window
	rl.InitWindow(screenWidth, screenHeight, "Ultimate Tic-Tac-Toe")
	rl.SetTargetFPS(60)

	// Main game loop
	for !rl.WindowShouldClose() {
		// Update
		updateGame(ui)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(boardBgColor)

		drawBoard(ui)
		drawPieces(ui)
		drawUI(ui)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

// updateGame handles game logic updates
func updateGame(ui *UIState) {
	// Update mouse position
	mousePosition := rl.GetMousePosition()

	// Check if mouse is over the board
	if mousePosition.X >= float32(ui.BoardOffset.X) &&
		mousePosition.X < float32(ui.BoardOffset.X+9*cellSize) &&
		mousePosition.Y >= float32(ui.BoardOffset.Y) &&
		mousePosition.Y < float32(ui.BoardOffset.Y+9*cellSize) {

		// Calculate which cell the mouse is hovering over
		cellX := int((mousePosition.X - float32(ui.BoardOffset.X)) / cellSize)
		cellY := int((mousePosition.Y - float32(ui.BoardOffset.Y)) / cellSize)

		ui.HoveredCell = Point{X: cellX, Y: cellY}

		// Calculate which small board and position this corresponds to
		boardRow := cellY / 3
		boardCol := cellX / 3
		boardIndex := boardRow*3 + boardCol

		localRow := cellY % 3
		localCol := cellX % 3
		localPosition := localRow*3 + localCol

		ui.HoveredBoard = boardIndex
		ui.HoveredPosition = localPosition
	} else {
		ui.HoveredCell = Point{X: -1, Y: -1}
		ui.HoveredBoard = -1
		ui.HoveredPosition = -1
	}

	// Handle mouse click
	if !ui.GameOver && rl.IsMouseButtonPressed(rl.MouseLeftButton) && ui.HoveredBoard != -1 {
		// Check if the move is valid
		validMove := false

		// Is move in active board, or any board if active board is -1
		if ui.Game.ActiveBoard == -1 || ui.Game.ActiveBoard == ui.HoveredBoard {
			// Convert to bit position to check if cell is empty
			bitPos := GetGlobalPosition(ui.HoveredBoard, ui.HoveredPosition)
			if ui.Game.IsEmpty(bitPos) {
				validMove = true
			}
		}

		if validMove {
			// Make the move
			ui.Game.MakeMove(ui.HoveredBoard, ui.HoveredPosition)

			// Store the last move for highlight
			ui.LastMove = ui.HoveredCell

			// Check if game is over
			if ui.Game.IsGameWon() {
				ui.GameOver = true
				ui.Winner = ui.Game.GetGameWinner()
			} else if ui.Game.IsGameDraw() {
				ui.GameOver = true
				ui.Winner = 0 // Draw
			}
		}
	}

	// New game on spacebar if game is over
	if ui.GameOver && rl.IsKeyPressed(rl.KeySpace) {
		ui.Game = NewGame()
		ui.GameOver = false
		ui.Winner = 0
		ui.LastMove = Point{X: -1, Y: -1}
	}
}

// drawBoard draws the game board
func drawBoard(ui *UIState) {
	offsetX := ui.BoardOffset.X
	offsetY := ui.BoardOffset.Y

	// Draw background for active board
	if ui.Game.ActiveBoard != -1 && !ui.GameOver {
		boardRow := ui.Game.ActiveBoard / 3
		boardCol := ui.Game.ActiveBoard % 3

		rl.DrawRectangle(
			int32(offsetX+boardCol*3*cellSize),
			int32(offsetY+boardRow*3*cellSize),
			int32(3*cellSize),
			int32(3*cellSize),
			activeBoardColor)
	}

	// Draw hover highlight
	if ui.HoveredCell.X != -1 && !ui.GameOver {
		// Check if the move would be valid
		validMove := false

		if ui.Game.ActiveBoard == -1 || ui.Game.ActiveBoard == ui.HoveredBoard {
			bitPos := GetGlobalPosition(ui.HoveredBoard, ui.HoveredPosition)
			if ui.Game.IsEmpty(bitPos) {
				validMove = true
			}
		}

		if validMove {
			rl.DrawRectangle(
				int32(offsetX+ui.HoveredCell.X*cellSize),
				int32(offsetY+ui.HoveredCell.Y*cellSize),
				int32(cellSize),
				int32(cellSize),
				hoverHighlightColor)
		}
	}

	// Draw won boards background
	for board := 0; board < 9; board++ {
		winner := ui.Game.GetBoardWinner(board)
		if winner > 0 {
			boardRow := board / 3
			boardCol := board % 3

			var color rl.Color
			if winner == 1 {
				color = wonBoardColorX
			} else {
				color = wonBoardColorO
			}

			rl.DrawRectangle(
				int32(offsetX+boardCol*3*cellSize),
				int32(offsetY+boardRow*3*cellSize),
				int32(3*cellSize),
				int32(3*cellSize),
				color)
		}
	}

	// Draw grid lines
	// Small grid lines
	for i := 0; i <= 9; i++ {
		// Vertical lines
		lineColor := gridLineColor
		if i%3 == 0 {
			lineColor = mainGridLineColor
		}

		rl.DrawLineEx(
			rl.NewVector2(float32(offsetX+i*cellSize), float32(offsetY)),
			rl.NewVector2(float32(offsetX+i*cellSize), float32(offsetY+9*cellSize)),
			2.0,
			lineColor)

		// Horizontal lines
		rl.DrawLineEx(
			rl.NewVector2(float32(offsetX), float32(offsetY+i*cellSize)),
			rl.NewVector2(float32(offsetX+9*cellSize), float32(offsetY+i*cellSize)),
			2.0,
			lineColor)
	}
}

// drawPieces draws all X and O pieces on the board
func drawPieces(ui *UIState) {
	offsetX := ui.BoardOffset.X
	offsetY := ui.BoardOffset.Y

	// Draw all pieces
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			bitPos := GetBitPosition(row, col)
			mask := uint64(1) << bitPos

			cellX := offsetX + col*cellSize
			cellY := offsetY + row*cellSize

			// Draw X
			if (ui.Game.XBoard & mask) != 0 {
				drawX(cellX, cellY, cellSize, xColor)
			}

			// Draw O
			if (ui.Game.OBoard & mask) != 0 {
				drawO(cellX, cellY, cellSize, oColor)
			}
		}
	}

	// Highlight last move
	if ui.LastMove.X != -1 {
		cellX := offsetX + ui.LastMove.X*cellSize
		cellY := offsetY + ui.LastMove.Y*cellSize

		rl.DrawRectangleLinesEx(
			rl.NewRectangle(
				float32(cellX),
				float32(cellY),
				float32(cellSize),
				float32(cellSize),
			),
			3.0,
			rl.Yellow)
	}
}

// drawX draws an X in a cell
func drawX(x, y, size int, color rl.Color) {
	padding := size / 4
	var thickness float32 = 3.0

	// Draw X (two diagonal lines)
	rl.DrawLineEx(
		rl.NewVector2(float32(x+padding), float32(y+padding)),
		rl.NewVector2(float32(x+size-padding), float32(y+size-padding)),
		thickness,
		color)

	rl.DrawLineEx(
		rl.NewVector2(float32(x+size-padding), float32(y+padding)),
		rl.NewVector2(float32(x+padding), float32(y+size-padding)),
		thickness,
		color)
}

// drawO draws an O in a cell
func drawO(x, y, size int, color rl.Color) {
	centerX := x + size/2
	centerY := y + size/2
	radius := (size / 2) - (size / 5)
	var thickness float32 = 3.0

	rl.DrawRingLines(
		rl.NewVector2(float32(centerX), float32(centerY)),
		float32(radius-int(float64(thickness/2))),
		float32(radius+int(float64(thickness/2))),
		0,
		360,
		0,
		color)
}

// drawUI draws game status and UI elements
func drawUI(ui *UIState) {
	// Draw status bar
	rl.DrawRectangle(0, 0, int32(screenWidth), 60, statusBgColor)

	// Draw game status
	var statusText string

	if ui.GameOver {
		if ui.Winner == 1 {
			statusText = "X wins! Press SPACE to play again."
		} else if ui.Winner == 2 {
			statusText = "O wins! Press SPACE to play again."
		} else {
			statusText = "Game is a draw! Press SPACE to play again."
		}
	} else {
		if ui.Game.XTurn {
			statusText = "X's turn"
		} else {
			statusText = "O's turn"
		}

		if ui.Game.ActiveBoard == -1 {
			statusText += " - Play in any board"
		} else {
			statusText += fmt.Sprintf(" - Play in board %d", ui.Game.ActiveBoard)
		}
	}

	rl.DrawText(
		statusText,
		screenWidth/2-rl.MeasureText(statusText, 24)/2,
		20,
		24,
		textColor)

	// Draw game instructions at the bottom
	instructionText := "Click on a valid cell to place your piece"
	rl.DrawText(
		instructionText,
		screenWidth/2-rl.MeasureText(instructionText, 20)/2,
		screenHeight-30,
		20,
		textColor)
}
