package game

import (
	"fmt"
	"sync"
)

// Engine manages the game logic and state
type Engine struct {
	games map[string]*GameState
	mutex sync.RWMutex
}

// NewEngine creates a new game engine
func NewEngine() *Engine {
	return &Engine{
		games: make(map[string]*GameState),
	}
}

// CreateGame creates a new game with the given ID
func (e *Engine) CreateGame(gameID string) *GameState {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	game := NewGame(gameID)
	e.games[gameID] = game
	return game
}

// GetGame retrieves a game by ID
func (e *Engine) GetGame(gameID string) (*GameState, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	game, exists := e.games[gameID]
	if !exists {
		return nil, fmt.Errorf("game with ID %s not found", gameID)
	}
	return game, nil
}

// MakeMove attempts to make a move on the board
func (e *Engine) MakeMove(gameID string, pos Position, player Player) (*GameState, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	game, exists := e.games[gameID]
	if !exists {
		return nil, fmt.Errorf("game with ID %s not found", gameID)
	}

	// Validate the move
	if err := e.validateMove(game, pos, player); err != nil {
		return nil, err
	}

	// Make the move
	game.Board.Set(pos, player)
	game.MoveCount++

	// Check for win condition
	if e.checkWin(game.Board, player) {
		game.Status = StatusWon
		game.Winner = player
	} else if game.Board.IsFull() {
		game.Status = StatusDraw
	} else {
		// Switch to next player
		game.CurrentPlayer = game.NextPlayer()
	}

	return game, nil
}

// validateMove checks if a move is valid
func (e *Engine) validateMove(game *GameState, pos Position, player Player) error {
	// Check if game is still ongoing
	if game.IsGameOver() {
		return fmt.Errorf("game is already over")
	}

	// Check if it's the correct player's turn
	if player != game.CurrentPlayer {
		return fmt.Errorf("it's not %s's turn", player)
	}

	// Check if position is valid
	if pos.Row < 0 || pos.Row > 2 || pos.Col < 0 || pos.Col > 2 {
		return fmt.Errorf("position %s is out of bounds", pos.String())
	}

	// Check if position is empty
	if !game.Board.IsEmpty(pos) {
		return fmt.Errorf("position %s is already occupied", pos.String())
	}

	return nil
}

// checkWin checks if the given player has won
func (e *Engine) checkWin(board Board, player Player) bool {
	// Check rows
	for row := 0; row < 3; row++ {
		if board[row][0] == player && board[row][1] == player && board[row][2] == player {
			return true
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		if board[0][col] == player && board[1][col] == player && board[2][col] == player {
			return true
		}
	}

	// Check diagonals
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true
	}

	return false
}

// ResetGame resets an existing game to initial state
func (e *Engine) ResetGame(gameID string) (*GameState, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	game, exists := e.games[gameID]
	if !exists {
		return nil, fmt.Errorf("game with ID %s not found", gameID)
	}

	// Reset to initial state
	game.Board = NewBoard()
	game.CurrentPlayer = PlayerX
	game.Status = StatusOngoing
	game.Winner = Empty
	game.MoveCount = 0

	return game, nil
}

// DeleteGame removes a game from the engine
func (e *Engine) DeleteGame(gameID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if _, exists := e.games[gameID]; !exists {
		return fmt.Errorf("game with ID %s not found", gameID)
	}

	delete(e.games, gameID)
	return nil
}

// ListGames returns all game IDs
func (e *Engine) ListGames() []string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	gameIDs := make([]string, 0, len(e.games))
	for id := range e.games {
		gameIDs = append(gameIDs, id)
	}
	return gameIDs
}

// GetAvailableMoves returns all valid positions for the current player
func (e *Engine) GetAvailableMoves(gameID string) ([]Position, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	game, exists := e.games[gameID]
	if !exists {
		return nil, fmt.Errorf("game with ID %s not found", gameID)
	}

	if game.IsGameOver() {
		return []Position{}, nil
	}

	var moves []Position
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			pos := Position{Row: row, Col: col}
			if game.Board.IsEmpty(pos) {
				moves = append(moves, pos)
			}
		}
	}

	return moves, nil
}

// AnalyzePosition provides analysis of the current game position
func (e *Engine) AnalyzePosition(gameID string) (string, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	game, exists := e.games[gameID]
	if !exists {
		return "", fmt.Errorf("game with ID %s not found", gameID)
	}

	if game.IsGameOver() {
		switch game.Status {
		case StatusWon:
			return fmt.Sprintf("Game over: %s wins!", game.Winner), nil
		case StatusDraw:
			return "Game over: It's a draw!", nil
		}
	}

	availableMoves, _ := e.GetAvailableMoves(gameID)
	return fmt.Sprintf("Current player: %s, Available moves: %d, Move count: %d",
		game.CurrentPlayer, len(availableMoves), game.MoveCount), nil
}
