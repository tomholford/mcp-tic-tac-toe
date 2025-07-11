package server

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"mcp-tic-tac-toe/game"
)

// TicTacToeServer wraps the game engine with MCP server functionality
type TicTacToeServer struct {
	mcpServer *server.MCPServer
	engine    *game.Engine
}

// NewTicTacToeServer creates a new MCP server for tic-tac-toe
func NewTicTacToeServer() *TicTacToeServer {
	s := &TicTacToeServer{
		engine: game.NewEngine(),
	}

	// Create MCP server with tool capabilities
	s.mcpServer = server.NewMCPServer(
		"Tic-Tac-Toe Game Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Register all tools
	s.registerTools()

	return s
}

// registerTools registers all MCP tools with their handlers
func (s *TicTacToeServer) registerTools() {
	// New game tool
	newGameTool := mcp.NewTool("new_game",
		mcp.WithDescription("Create a new tic-tac-toe game"),
		mcp.WithString("game_id",
			mcp.Description("Optional game ID. If not provided, a random ID will be generated"),
		),
	)
	s.mcpServer.AddTool(newGameTool, s.handleNewGame)

	// Make move tool
	makeMoveTool := mcp.NewTool("make_move",
		mcp.WithDescription("Make a move on the tic-tac-toe board"),
		mcp.WithString("game_id",
			mcp.Required(),
			mcp.Description("ID of the game to make a move in"),
		),
		mcp.WithString("position",
			mcp.Required(),
			mcp.Description("Position to place mark (A1-C3 format)"),
			mcp.Enum("A1", "A2", "A3", "B1", "B2", "B3", "C1", "C2", "C3"),
		),
		mcp.WithString("player",
			mcp.Required(),
			mcp.Enum("X", "O"),
			mcp.Description("Player making the move"),
		),
	)
	s.mcpServer.AddTool(makeMoveTool, s.handleMakeMove)

	// Get board tool
	getBoardTool := mcp.NewTool("get_board",
		mcp.WithDescription("Get the current board state"),
		mcp.WithString("game_id",
			mcp.Required(),
			mcp.Description("ID of the game to get board state for"),
		),
	)
	s.mcpServer.AddTool(getBoardTool, s.handleGetBoard)

	// Get game status tool
	getStatusTool := mcp.NewTool("get_status",
		mcp.WithDescription("Get the current game status and winner"),
		mcp.WithString("game_id",
			mcp.Required(),
			mcp.Description("ID of the game to get status for"),
		),
	)
	s.mcpServer.AddTool(getStatusTool, s.handleGetStatus)

	// Reset game tool
	resetGameTool := mcp.NewTool("reset_game",
		mcp.WithDescription("Reset an existing game to initial state"),
		mcp.WithString("game_id",
			mcp.Required(),
			mcp.Description("ID of the game to reset"),
		),
	)
	s.mcpServer.AddTool(resetGameTool, s.handleResetGame)

	// Get available moves tool
	getAvailableMovesTool := mcp.NewTool("get_available_moves",
		mcp.WithDescription("Get all available/valid moves for current player"),
		mcp.WithString("game_id",
			mcp.Required(),
			mcp.Description("ID of the game to get available moves for"),
		),
	)
	s.mcpServer.AddTool(getAvailableMovesTool, s.handleGetAvailableMoves)

	// Analyze position tool
	analyzePositionTool := mcp.NewTool("analyze_position",
		mcp.WithDescription("Analyze the current game position"),
		mcp.WithString("game_id",
			mcp.Required(),
			mcp.Description("ID of the game to analyze"),
		),
	)
	s.mcpServer.AddTool(analyzePositionTool, s.handleAnalyzePosition)

	// List games tool
	listGamesTool := mcp.NewTool("list_games",
		mcp.WithDescription("List all active game IDs"),
	)
	s.mcpServer.AddTool(listGamesTool, s.handleListGames)
}

// generateGameID creates a random game ID
func generateGameID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return "game-" + hex.EncodeToString(bytes)
}

// GetMCPServer returns the underlying MCP server
func (s *TicTacToeServer) GetMCPServer() *server.MCPServer {
	return s.mcpServer
}

// ServeStdio starts the server using stdio transport
func (s *TicTacToeServer) ServeStdio() error {
	log.Println("Starting Tic-Tac-Toe MCP server with stdio transport...")
	return server.ServeStdio(s.mcpServer)
}

// ServeSSE starts the server using Server-Sent Events transport
func (s *TicTacToeServer) ServeSSE(addr string) error {
	log.Printf("Starting Tic-Tac-Toe MCP server with SSE transport on %s...", addr)
	sseServer := server.NewSSEServer(s.mcpServer)
	return sseServer.Start(addr)
}

// ServeStreamableHTTP starts the server using Streamable HTTP transport
func (s *TicTacToeServer) ServeStreamableHTTP(addr string) error {
	log.Printf("Starting Tic-Tac-Toe MCP server with Streamable HTTP transport on %s...", addr)
	httpServer := server.NewStreamableHTTPServer(s.mcpServer)
	return httpServer.Start(addr)
}
