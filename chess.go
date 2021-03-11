package main

import (
	"fmt"
	"strings"
	"sync"
)

var (
	player1, player2 Player
	mu               = &sync.Mutex{}
	// We use the map to initialize the pieces, and later on update the piece slice according to the pieces the player actually has
	pieceMap = map[string]int{"king": 1, "queen": 1, "rook": 2, "bishop": 2, "knight": 2, "pawn": 8}
	// sessionPlayers holdes the current session Players
	sessionPlayers []Player
)

/*
Player interface containes all the methods of a player in the game.
This interface represents the Player Abstract Class on your code
*/
type Player interface {
	GetName() string
	GetColor() string
	GetPieces() []piece
	GetAvailability() bool
	CreatePieces(string, int) bool
	MovePiece(string, string) error
	KillKing() bool
}

// chessPlayer contains each player attributes
type chessPlayer struct {
	name           string
	color          string
	playablePieces []piece
}

// piece contains each piece attributes
type piece struct {
	name            string
	color           string
	currentPosition string
	exist           bool
	allowMoves      []string
}

// MovePiece lock the state, and allows a player to move a piece on his turn. returns err if game ended.
func (v chessPlayer) MovePiece(pName, location string) error {
	defer mu.Unlock() // Release the state at the end of the func
	mu.Lock()         /* In order to keep the game in sync (i.e - a player cannot move a piece until the 2'nd player is done),
	we are locking the game */

	// First, we want to make sure if the game is not already ended
	if IsGameEnded() {
		return fmt.Errorf("Sorry, the game already ended. The winner is: %s", sessionPlayers[0].GetName())
	}
	fmt.Printf("Player %s is moving %s to %s.\n", v.name, strings.Title(pName), location)
	if msg := gameCoordinator(v, strings.ToLower(pName), location); msg != "" { // gameCoordinator is the main Chess game logic function - see doc
		fmt.Printf("%s cannot move to %s. Move is illegal.\n", strings.Title(pName), location)
		return nil
	}
	fmt.Printf("%s moved to %s.\n", strings.Title(pName), location)
	return nil
}

// KillKing remove the kind from the board, just to flag a game end
func (v *chessPlayer) KillKing() bool {
	for i, p := range v.playablePieces {
		if strings.ToLower(p.name) == "king" {
			v.playablePieces[i] = v.playablePieces[len(v.playablePieces)-1]
			v.playablePieces = v.playablePieces[:len(v.playablePieces)-1]
		}
	}
	// Also, remove the player from the session
	for i, plr := range sessionPlayers {
		if plr.GetName() == v.name {
			sessionPlayers[i] = sessionPlayers[len(sessionPlayers)-1]
			sessionPlayers = sessionPlayers[:len(sessionPlayers)-1]
		}
	}
	return true
}

// GetName return the player's name
func (v chessPlayer) GetName() string {
	return v.name
}

// GetColor return the color of the player
func (v chessPlayer) GetColor() string {
	return v.color
}

// GetPieces returns which pieces left
func (v chessPlayer) GetPieces() []piece {
	return v.playablePieces
}

// GetAvailability returns if player have lost the game
func (v chessPlayer) GetAvailability() bool {
	for _, p := range v.playablePieces {
		if strings.ToLower(p.name) == "king" { // gameCoordinator() make sure to remove the King when a player have lost
			return true
		}
	}
	return false
}

// CreatePieces populates the piece struct
func (v *chessPlayer) CreatePieces(pName string, pCount int) bool {
	// Create as many piece as needed based on role (King, Pawn, etc)
	for i := 0; i < pCount; i++ {
		v.playablePieces = append(v.playablePieces, piece{
			color:           v.color,
			name:            strings.ToLower(pName),
			currentPosition: "a1", // TODO: placementLogic(), a function that will place the piece in the game starting position
			exist:           true,
			allowMoves:      []string{"blabla"}, // TODO: initPieceMoves(), a function that will set on what direction a piece can move to
		})
	}
	return true
}

// IsGameEnded return a bool on the game status. This funcion can be wrapped using a for loop.
func IsGameEnded() bool {
	if !player1.GetAvailability() || !player2.GetAvailability() {
		return true
	}
	return false
}

func gameCoordinator(player chessPlayer, pName, location string) string {
	/*
		TODO:
		gameCoordinator will coordinate the game according to Chess rules.
		Some logic I decided not to implaement in this task contains:
		1. Get the piece move coordinations (pName + location), and act upon it:
			a. Validate that the piece type (`pName`) can move as requested (based on `piece.allowMoves`)
			b. Validate that the path is empty - if so, move the desired location and update `piece.currentPosition`
			c. If the square is not empty, captured the target piece and "remove it from the table" (switch `piece.exist` to false and remove it from `playablePieces` slice).
		2. Validate and catch the different Chess ending (for example, the King is under threat and has no way to move),
		   and update piece.exist` and `playablePieces` accordingly.
		3. Return the outcome of the move ("The queen is on A8, and the Pawn has been removed from the board.")

	*/
	return ""
}

func placePieces([]Player) {
	for _, plr := range sessionPlayers {
		for pieceName, pieceCount := range pieceMap {
			/*
				As requested I'm not actually going to use those vars to place each one of the pieces in place
				based on Chess rules, but I would use them for placement and for creating the piece object out of it,
				All pieces will be appended to the playablePieces slice which holds all the pieces a player has.
			*/
			plr.CreatePieces(pieceName, pieceCount)
		}
		fmt.Printf("Player %s is ready on the board. Going to play %s\n", plr.GetName(), plr.GetColor())
	}
}

/*
The init function will be executed once, no matter how many times the package is imported,
to make sure we have only 2 players.
*/
func init() {
	player1 = &chessPlayer{"Garry Kasparov", "White", []piece{}}
	player2 = &chessPlayer{"Deep Blue", "Black", []piece{}}
	sessionPlayers = append(sessionPlayers, player1, player2)
}

/*
This method starts the chess game, and returns the winner.
*/
func main() {
	fmt.Printf("Let's play some Chess!\n♔ ♕ ♖ ♗ ♘ ♙ ♚ ♛ ♜ ♝ ♞ ♟\n")
	// initialize the board and place the pieces
	placePieces(sessionPlayers)

	//To play the game, you can import the package, and wrapped in a loop around `IsGameEnded()` -
	for !IsGameEnded() {
		player1.MovePiece("queen", "f6")
		player2.MovePiece("bishop", "a8")
		player1.MovePiece("pawn", "h4")
		player1.KillKing()              // Let's make player1 lose the game
		player1.MovePiece("king", "h6") // Won't be evaluated
	}
	fmt.Printf("The game has ended!\nThe winner is: %s\n", sessionPlayers[0].GetName())

	/* Another option, is to check each one of the moves (e.g - send all moves to channel and recive one by one to MovePiece(), etc).
	   The function will return an error if the game already ended */
	placePieces(sessionPlayers)
	if err := player1.MovePiece("queen", "f6"); err != nil {
		fmt.Println(err)
	}
	/*
		Some more moves
	*/
	player1.KillKing()
	if err := player1.MovePiece("bishop", "a8"); err != nil {
		fmt.Println(err)
	}
}
