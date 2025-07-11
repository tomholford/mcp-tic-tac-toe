package game

import (
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}
	if len(engine.games) != 0 {
		t.Error("New engine should have no games")
	}
}

func TestCreateGame(t *testing.T) {
	engine := NewEngine()
	game := engine.CreateGame("test-game")

	if game == nil {
		t.Fatal("CreateGame() returned nil")
	}
	if game.GameID != "test-game" {
		t.Errorf("Expected game ID 'test-game', got '%s'", game.GameID)
	}
	if game.CurrentPlayer != PlayerX {
		t.Errorf("Expected current player to be X, got %s", game.CurrentPlayer)
	}
	if game.Status != StatusOngoing {
		t.Errorf("Expected status to be ongoing, got %s", game.Status)
	}
}

func TestMakeMove(t *testing.T) {
	engine := NewEngine()
	engine.CreateGame("test-game")

	// Test valid move
	pos, _ := ParsePosition("A1")
	updatedGame, err := engine.MakeMove("test-game", pos, PlayerX)
	if err != nil {
		t.Fatalf("MakeMove() failed: %v", err)
	}

	if updatedGame.Board.Get(pos) != PlayerX {
		t.Error("Move was not applied to board")
	}
	if updatedGame.CurrentPlayer != PlayerO {
		t.Error("Player should switch after valid move")
	}
	if updatedGame.MoveCount != 1 {
		t.Errorf("Expected move count 1, got %d", updatedGame.MoveCount)
	}
}

func TestInvalidMoves(t *testing.T) {
	engine := NewEngine()
	engine.CreateGame("test-game")

	// Test wrong player
	pos, _ := ParsePosition("A1")
	_, err := engine.MakeMove("test-game", pos, PlayerO)
	if err == nil {
		t.Error("Should reject wrong player's move")
	}

	// Make valid move first
	engine.MakeMove("test-game", pos, PlayerX)

	// Test occupied position
	_, err = engine.MakeMove("test-game", pos, PlayerO)
	if err == nil {
		t.Error("Should reject move to occupied position")
	}
}

func TestWinConditions(t *testing.T) {
	engine := NewEngine()

	// Test row win
	engine.CreateGame("row-win")
	positions := []string{"A1", "A2", "B1", "B2", "C1"} // X wins with row 1
	players := []Player{PlayerX, PlayerO, PlayerX, PlayerO, PlayerX}

	for i, posStr := range positions {
		pos, _ := ParsePosition(posStr)
		updatedGame, err := engine.MakeMove("row-win", pos, players[i])
		if err != nil {
			t.Fatalf("Move %d failed: %v", i, err)
		}

		if i == 4 { // Last move should trigger win
			if updatedGame.Status != StatusWon {
				t.Error("Game should be won after row completion")
			}
			if updatedGame.Winner != PlayerX {
				t.Error("PlayerX should be the winner")
			}
		}
	}
}

func TestDraw(t *testing.T) {
	engine := NewEngine()
	engine.CreateGame("draw-game")

	// Create a draw scenario - arrange moves to avoid any wins
	// Final board should look like:
	// X O X
	// O O X
	// X X O
	moves := []struct {
		pos    string
		player Player
	}{
		{"A1", PlayerX}, {"A2", PlayerO}, {"A3", PlayerX}, // Row 1: X O X
		{"B1", PlayerO}, {"C1", PlayerX}, {"B2", PlayerO}, // O in B1, X in C1, O in B2
		{"B3", PlayerX}, {"C3", PlayerO}, {"C2", PlayerX}, // Complete the draw
	}

	for i, move := range moves {
		pos, _ := ParsePosition(move.pos)
		updatedGame, err := engine.MakeMove("draw-game", pos, move.player)
		if err != nil {
			t.Fatalf("Move %d (%s by %s) failed: %v", i, move.pos, move.player, err)
		}

		if i == 8 { // Last move should trigger draw
			if updatedGame.Status != StatusDraw {
				t.Errorf("Game should be a draw, but status is %s", updatedGame.Status)
			}
		}
	}
}

func TestGetAvailableMoves(t *testing.T) {
	engine := NewEngine()
	engine.CreateGame("moves-test")

	// Initially all 9 positions should be available
	moves, err := engine.GetAvailableMoves("moves-test")
	if err != nil {
		t.Fatalf("GetAvailableMoves() failed: %v", err)
	}
	if len(moves) != 9 {
		t.Errorf("Expected 9 available moves, got %d", len(moves))
	}

	// Make a move
	pos, _ := ParsePosition("A1")
	engine.MakeMove("moves-test", pos, PlayerX)

	// Should have 8 moves left
	moves, _ = engine.GetAvailableMoves("moves-test")
	if len(moves) != 8 {
		t.Errorf("Expected 8 available moves after one move, got %d", len(moves))
	}
}

func TestPositionParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected Position
		hasError bool
	}{
		{"A1", Position{0, 0}, false},
		{"B2", Position{1, 1}, false},
		{"C3", Position{2, 2}, false},
		{"D1", Position{}, true}, // Invalid column
		{"A4", Position{}, true}, // Invalid row
		{"", Position{}, true},   // Empty string
		{"A", Position{}, true},  // Too short
	}

	for _, tc := range testCases {
		pos, err := ParsePosition(tc.input)
		if tc.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s'", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
			}
			if pos != tc.expected {
				t.Errorf("For input '%s', expected %+v, got %+v", tc.input, tc.expected, pos)
			}
		}
	}
}
