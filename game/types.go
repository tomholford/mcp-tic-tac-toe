package game

import "fmt"

// Player represents a player in the game
type Player string

const (
	PlayerX Player = "X"
	PlayerO Player = "O"
	Empty   Player = ""
)

// Position represents a position on the board using chess-like notation
type Position struct {
	Row int // 0-2
	Col int // 0-2
}

// String returns the position in A1-C3 format
func (p Position) String() string {
	if p.Row < 0 || p.Row > 2 || p.Col < 0 || p.Col > 2 {
		return "Invalid"
	}
	return fmt.Sprintf("%c%d", 'A'+p.Col, p.Row+1)
}

// ParsePosition converts A1-C3 notation to Position
func ParsePosition(pos string) (Position, error) {
	if len(pos) != 2 {
		return Position{}, fmt.Errorf("position must be 2 characters (e.g., A1)")
	}

	col := int(pos[0] - 'A')
	row := int(pos[1] - '1')

	if col < 0 || col > 2 || row < 0 || row > 2 {
		return Position{}, fmt.Errorf("position must be A1-C3")
	}

	return Position{Row: row, Col: col}, nil
}

// GameStatus represents the current state of the game
type GameStatus string

const (
	StatusOngoing GameStatus = "ongoing"
	StatusWon     GameStatus = "won"
	StatusDraw    GameStatus = "draw"
)

// Board represents the 3x3 tic-tac-toe board
type Board [3][3]Player

// NewBoard creates a new empty board
func NewBoard() Board {
	return Board{}
}

// Get returns the player at the given position
func (b *Board) Get(pos Position) Player {
	return b[pos.Row][pos.Col]
}

// Set places a player at the given position
func (b *Board) Set(pos Position, player Player) {
	b[pos.Row][pos.Col] = player
}

// IsEmpty checks if a position is empty
func (b *Board) IsEmpty(pos Position) bool {
	return b[pos.Row][pos.Col] == Empty
}

// IsFull checks if the board is completely filled
func (b *Board) IsFull() bool {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if b[row][col] == Empty {
				return false
			}
		}
	}
	return true
}

// String returns a formatted string representation of the board
func (b *Board) String() string {
	result := "  A B C\n"
	for row := 0; row < 3; row++ {
		result += fmt.Sprintf("%d ", row+1)
		for col := 0; col < 3; col++ {
			cell := string(b[row][col])
			if cell == "" {
				cell = "·"
			}
			result += cell
			if col < 2 {
				result += " "
			}
		}
		result += "\n"
	}
	return result
}

// GameState represents the complete state of a game
type GameState struct {
	Board         Board
	CurrentPlayer Player
	Status        GameStatus
	Winner        Player
	MoveCount     int
	GameID        string
}

// NewGame creates a new game state
func NewGame(gameID string) *GameState {
	return &GameState{
		Board:         NewBoard(),
		CurrentPlayer: PlayerX, // X always goes first
		Status:        StatusOngoing,
		Winner:        Empty,
		MoveCount:     0,
		GameID:        gameID,
	}
}

// NextPlayer returns the next player to move
func (g *GameState) NextPlayer() Player {
	if g.CurrentPlayer == PlayerX {
		return PlayerO
	}
	return PlayerX
}

// IsGameOver checks if the game has ended
func (g *GameState) IsGameOver() bool {
	return g.Status != StatusOngoing
}
