package main

import (
	"fmt"

	"github.com/not-swym/u3t/game" // Import our game package
)

func main() {
	fmt.Println("Starting Ultimate Tic-Tac-Toe...")

	// Create new game instance
	gameState := game.NewGame()

	// Initialize and run the UI with our game state
	ui := game.NewUIState(gameState)
	game.RunGame(ui)
}
