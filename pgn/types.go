package pgn

// Game is a structure representing a complete chess game containing metadata (Tags) and the actual
// moves that made up the chess game (Moves).
type Game struct {
	// Tags are any unstructure metadata belonging to a chess game.
	Tags map[string]string
	// Moves are the moves and annotaitons that make up a chess game
	Moves []Move
}

// Move is a structure that defines a chess move
type Move struct {
	// Number is the index of the move within a game
	Number int32
	// Move is the [algebraic notation](https://en.wikipedia.org/wiki/Algebraic_notation_(chess)) represntation of the move in SAN format
	Move string
	// Annotation is a comment assigned to the move
	Annotation string
	// Nag is a Numeric Annotation Glyph (ie. !! or !? or one of those crazy chess characters)
	Nag string
	// Alternatices is a list of alternate moves and refutations that could have been played instead of this move
	Alternatives []Move
}
