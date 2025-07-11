package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"mcp-tic-tac-toe/game"
)

// handleNewGame creates a new game
func (s *TicTacToeServer) handleNewGame(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	// Get game ID or generate one
	var gameID string
	if id, ok := arguments["game_id"].(string); ok && id != "" {
		gameID = id
	} else {
		gameID = generateGameID()
	}

	// Create the game
	gameState := s.engine.CreateGame(gameID)

	response := fmt.Sprintf("New game created with ID: %s\nStarting player: %s\nInitial board:\n%s",
		gameState.GameID, gameState.CurrentPlayer, gameState.Board.String())

	return mcp.NewToolResultText(response), nil
}

// handleMakeMove processes a move request
func (s *TicTacToeServer) handleMakeMove(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameters
	gameID, err := request.RequireString("game_id")
	if err != nil {
		return mcp.NewToolResultError("game_id is required"), nil
	}

	positionStr, err := request.RequireString("position")
	if err != nil {
		return mcp.NewToolResultError("position is required"), nil
	}

	playerStr, err := request.RequireString("player")
	if err != nil {
		return mcp.NewToolResultError("player is required"), nil
	}

	// Parse position
	position, err := game.ParsePosition(positionStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid position: %v", err)), nil
	}

	// Parse player
	var player game.Player
	switch strings.ToUpper(playerStr) {
	case "X":
		player = game.PlayerX
	case "O":
		player = game.PlayerO
	default:
		return mcp.NewToolResultError("Player must be 'X' or 'O'"), nil
	}

	// Make the move
	gameState, err := s.engine.MakeMove(gameID, position, player)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Move failed: %v", err)), nil
	}

	// Build response
	response := fmt.Sprintf("Move successful: %s placed %s at %s\n\nUpdated board:\n%s",
		player, player, positionStr, gameState.Board.String())

	// Add game status
	if gameState.IsGameOver() {
		switch gameState.Status {
		case game.StatusWon:
			response += fmt.Sprintf("\n🎉 Game Over! %s wins!", gameState.Winner)
		case game.StatusDraw:
			response += "\nGame Over! It's a draw!"
		}
	} else {
		response += fmt.Sprintf("\nNext player: %s", gameState.CurrentPlayer)
	}

	return mcp.NewToolResultText(response), nil
}

// handleGetBoard returns the current board state
func (s *TicTacToeServer) handleGetBoard(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gameID, err := request.RequireString("game_id")
	if err != nil {
		return mcp.NewToolResultError("game_id is required"), nil
	}

	gameState, err := s.engine.GetGame(gameID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Game not found: %v", err)), nil
	}

	response := fmt.Sprintf("Game ID: %s\nCurrent board:\n%s\nCurrent player: %s\nMove count: %d",
		gameState.GameID, gameState.Board.String(), gameState.CurrentPlayer, gameState.MoveCount)

	return mcp.NewToolResultText(response), nil
}

// handleGetStatus returns the current game status
func (s *TicTacToeServer) handleGetStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gameID, err := request.RequireString("game_id")
	if err != nil {
		return mcp.NewToolResultError("game_id is required"), nil
	}

	gameState, err := s.engine.GetGame(gameID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Game not found: %v", err)), nil
	}

	var response string
	switch gameState.Status {
	case game.StatusOngoing:
		response = fmt.Sprintf("Game Status: Ongoing\nCurrent player: %s\nMove count: %d",
			gameState.CurrentPlayer, gameState.MoveCount)
	case game.StatusWon:
		response = fmt.Sprintf("Game Status: Completed\nWinner: %s\nTotal moves: %d",
			gameState.Winner, gameState.MoveCount)
	case game.StatusDraw:
		response = fmt.Sprintf("Game Status: Draw\nTotal moves: %d", gameState.MoveCount)
	}

	return mcp.NewToolResultText(response), nil
}

// handleResetGame resets a game to initial state
func (s *TicTacToeServer) handleResetGame(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gameID, err := request.RequireString("game_id")
	if err != nil {
		return mcp.NewToolResultError("game_id is required"), nil
	}

	gameState, err := s.engine.ResetGame(gameID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Reset failed: %v", err)), nil
	}

	response := fmt.Sprintf("Game %s has been reset\nStarting player: %s\nBoard:\n%s",
		gameState.GameID, gameState.CurrentPlayer, gameState.Board.String())

	return mcp.NewToolResultText(response), nil
}

// handleGetAvailableMoves returns all valid moves for the current player
func (s *TicTacToeServer) handleGetAvailableMoves(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gameID, err := request.RequireString("game_id")
	if err != nil {
		return mcp.NewToolResultError("game_id is required"), nil
	}

	moves, err := s.engine.GetAvailableMoves(gameID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get moves: %v", err)), nil
	}

	if len(moves) == 0 {
		return mcp.NewToolResultText("No available moves (game is over)"), nil
	}

	// Convert positions to strings
	moveStrs := make([]string, len(moves))
	for i, pos := range moves {
		moveStrs[i] = pos.String()
	}

	response := fmt.Sprintf("Available moves (%d): %s", len(moves), strings.Join(moveStrs, ", "))
	return mcp.NewToolResultText(response), nil
}

// handleAnalyzePosition provides analysis of the current position
func (s *TicTacToeServer) handleAnalyzePosition(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gameID, err := request.RequireString("game_id")
	if err != nil {
		return mcp.NewToolResultError("game_id is required"), nil
	}

	analysis, err := s.engine.AnalyzePosition(gameID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	// Get additional details for enhanced analysis
	gameState, err := s.engine.GetGame(gameID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Game not found: %v", err)), nil
	}

	response := fmt.Sprintf("Position Analysis for Game %s:\n%s\n\nCurrent board:\n%s",
		gameID, analysis, gameState.Board.String())

	return mcp.NewToolResultText(response), nil
}

// handleListGames returns all active game IDs
func (s *TicTacToeServer) handleListGames(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gameIDs := s.engine.ListGames()

	if len(gameIDs) == 0 {
		return mcp.NewToolResultText("No active games"), nil
	}

	response := fmt.Sprintf("Active games (%d): %s", len(gameIDs), strings.Join(gameIDs, ", "))
	return mcp.NewToolResultText(response), nil
}
