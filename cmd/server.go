package main

import (
	"flag"
	"log"
	"os"

	"mcp-tic-tac-toe/server"
)

func main() {
	var transport = flag.String("transport", "stdio", "Transport method: stdio, sse, or http")
	var addr = flag.String("addr", ":8080", "Address to listen on (for sse/http transport)")
	flag.Parse()

	// Create the tic-tac-toe MCP server
	gameServer := server.NewTicTacToeServer()

	log.Printf("Starting MCP Tic-Tac-Toe server with %s transport", *transport)

	var err error
	switch *transport {
	case "stdio":
		err = gameServer.ServeStdio()
	case "sse":
		err = gameServer.ServeSSE(*addr)
	case "http":
		err = gameServer.ServeStreamableHTTP(*addr)
	default:
		log.Fatalf("Unknown transport: %s (supported: stdio, sse, http)", *transport)
	}

	if err != nil {
		log.Fatalf("Server failed: %v", err)
		os.Exit(1)
	}
}
