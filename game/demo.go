package game

import "fmt"

// Demo demonstrates the core game functionality
func Demo() {
	fmt.Println("=== Tic-Tac-Toe Core Game Engine Demo ===")

	// Create a new game engine
	engine := NewEngine()

	// Create a new game
	game := engine.CreateGame("demo-game")
	fmt.Printf("Created new game: %s\n", game.GameID)
	fmt.Printf("Starting player: %s\n", game.CurrentPlayer)
	fmt.Printf("Initial board:\n%s\n", game.Board.String())

	// Simulate a game
	moves := []struct {
		pos    string
		player Player
		desc   string
	}{
		{"B2", PlayerX, "X takes center"},
		{"A1", PlayerO, "O takes corner"},
		{"A2", PlayerX, "X blocks O's row"},
		{"C2", PlayerO, "O takes middle right"},
		{"C3", PlayerX, "X takes bottom right"},
		{"A3", PlayerO, "O completes first column - O wins!"},
	}

	for i, move := range moves {
		fmt.Printf("Move %d: %s\n", i+1, move.desc)

		pos, err := ParsePosition(move.pos)
		if err != nil {
			fmt.Printf("Error parsing position: %v\n", err)
			continue
		}

		updatedGame, err := engine.MakeMove("demo-game", pos, move.player)
		if err != nil {
			fmt.Printf("Error making move: %v\n", err)
			continue
		}

		fmt.Printf("Board after %s plays %s:\n%s", move.player, move.pos, updatedGame.Board.String())

		if updatedGame.IsGameOver() {
			switch updatedGame.Status {
			case StatusWon:
				fmt.Printf("🎉 Game Over! %s wins!\n\n", updatedGame.Winner)
			case StatusDraw:
				fmt.Printf("Game Over! It's a draw!\n\n")
			}
			break
		} else {
			fmt.Printf("Next player: %s\n\n", updatedGame.CurrentPlayer)
		}
	}

	// Show available moves functionality
	fmt.Println("=== Demonstrating Available Moves ===")
	engine.CreateGame("moves-demo")

	// Make a few moves
	pos1, _ := ParsePosition("A1")
	updatedGame, _ := engine.MakeMove("moves-demo", pos1, PlayerX)
	pos2, _ := ParsePosition("B2")
	updatedGame, _ = engine.MakeMove("moves-demo", pos2, PlayerO)

	fmt.Printf("Board with 2 moves:\n%s", updatedGame.Board.String())

	availableMoves, _ := engine.GetAvailableMoves("moves-demo")
	fmt.Printf("Available moves (%d): ", len(availableMoves))
	for i, pos := range availableMoves {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(pos.String())
	}
	fmt.Println()

	// Show analysis
	analysis, _ := engine.AnalyzePosition("moves-demo")
	fmt.Printf("Position analysis: %s\n", analysis)
}
