package server

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewTicTacToeServer(t *testing.T) {
	server := NewTicTacToeServer()
	if server == nil {
		t.Fatal("NewTicTacToeServer() returned nil")
	}
	if server.mcpServer == nil {
		t.Fatal("MCP server not initialized")
	}
	if server.engine == nil {
		t.Fatal("Game engine not initialized")
	}
}

func TestNewGameTool(t *testing.T) {
	server := NewTicTacToeServer()
	ctx := context.Background()

	// Test with auto-generated ID
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "new_game",
			Arguments: map[string]interface{}{},
		},
	}

	result, err := server.handleNewGame(ctx, request)
	if err != nil {
		t.Fatalf("handleNewGame failed: %v", err)
	}

	if result == nil {
		t.Fatal("handleNewGame returned nil result")
	}

	// Check that response contains expected information
	response := getTextFromResult(result)
	if !strings.Contains(response, "New game created") {
		t.Error("Response should contain 'New game created'")
	}
	if !strings.Contains(response, "game-") {
		t.Error("Response should contain generated game ID")
	}
	if !strings.Contains(response, "Starting player: X") {
		t.Error("Response should show starting player")
	}

	// Test with specific ID
	request.Params.Arguments = map[string]interface{}{
		"game_id": "test-game-123",
	}

	result, err = server.handleNewGame(ctx, request)
	if err != nil {
		t.Fatalf("handleNewGame with specific ID failed: %v", err)
	}

	response = getTextFromResult(result)
	if !strings.Contains(response, "test-game-123") {
		t.Error("Response should contain specified game ID")
	}
}

func TestMakeMoveTool(t *testing.T) {
	server := NewTicTacToeServer()
	ctx := context.Background()

	// Create a game first
	server.engine.CreateGame("test-game")

	// Test valid move
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "make_move",
			Arguments: map[string]interface{}{
				"game_id":  "test-game",
				"position": "A1",
				"player":   "X",
			},
		},
	}

	result, err := server.handleMakeMove(ctx, request)
	if err != nil {
		t.Fatalf("handleMakeMove failed: %v", err)
	}

	response := getTextFromResult(result)
	if !strings.Contains(response, "Move successful") {
		t.Error("Response should indicate move success")
	}
	if !strings.Contains(response, "Next player: O") {
		t.Error("Response should show next player")
	}

	// Test invalid move (wrong player's turn)
	request.Params.Arguments = map[string]interface{}{
		"game_id":  "test-game",
		"position": "A2",
		"player":   "X", // Should be O's turn now
	}
	result, err = server.handleMakeMove(ctx, request)
	if err != nil {
		t.Fatalf("handleMakeMove should return error result, not error: %v", err)
	}

	if !result.IsError {
		t.Error("Result should indicate error for wrong player")
	}
}

func TestGetBoardTool(t *testing.T) {
	server := NewTicTacToeServer()
	ctx := context.Background()

	// Create a game
	server.engine.CreateGame("test-game")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_board",
			Arguments: map[string]interface{}{
				"game_id": "test-game",
			},
		},
	}

	result, err := server.handleGetBoard(ctx, request)
	if err != nil {
		t.Fatalf("handleGetBoard failed: %v", err)
	}

	response := getTextFromResult(result)
	if !strings.Contains(response, "Current board:") {
		t.Error("Response should contain 'Current board:'")
	}
	if !strings.Contains(response, "Current player: X") {
		t.Error("Response should show current player")
	}
	if !strings.Contains(response, "Move count: 0") {
		t.Error("Response should show move count")
	}

	// Test non-existent game
	request.Params.Arguments = map[string]interface{}{
		"game_id": "non-existent",
	}
	result, err = server.handleGetBoard(ctx, request)
	if err != nil {
		t.Fatalf("handleGetBoard should return error result, not error: %v", err)
	}

	if !result.IsError {
		t.Error("Result should indicate error for non-existent game")
	}
}

func TestListGamesTool(t *testing.T) {
	server := NewTicTacToeServer()
	ctx := context.Background()

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "list_games",
			Arguments: map[string]interface{}{},
		},
	}

	// Test with no games
	result, err := server.handleListGames(ctx, request)
	if err != nil {
		t.Fatalf("handleListGames failed: %v", err)
	}

	response := getTextFromResult(result)
	if !strings.Contains(response, "No active games") {
		t.Error("Response should indicate no active games")
	}

	// Create some games
	server.engine.CreateGame("game1")
	server.engine.CreateGame("game2")

	result, err = server.handleListGames(ctx, request)
	if err != nil {
		t.Fatalf("handleListGames failed with games: %v", err)
	}

	response = getTextFromResult(result)
	if !strings.Contains(response, "Active games (2):") {
		t.Error("Response should show 2 active games")
	}
	if !strings.Contains(response, "game1") || !strings.Contains(response, "game2") {
		t.Error("Response should contain both game IDs")
	}
}

// Helper functions
func getTextFromResult(result *mcp.CallToolResult) string {
	if len(result.Content) == 0 {
		return ""
	}
	if textContent, ok := result.Content[0].(mcp.TextContent); ok {
		return textContent.Text
	}
	return ""
}
