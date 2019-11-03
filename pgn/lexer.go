package pgn

import (
	"io"
	"text/scanner"
)

// Tok defines the types of legal tokens that may be used in a pgn
type Tok int

const eof = rune(-1)
const (
	// Illegal token is not part of the pgn spec
	Illegal Tok = iota
	// EOF token is the end of the input
	EOF
	// Ws token is a white space token and is omitted from most
	// of the scanner outputs
	Ws

	// LBrace is a '[' symbol
	LBrace
	// RBrace is a ']' symbol
	RBrace
	// LParen is a '(' symbol
	LParen
	// RParen is a ')' symbol
	RParen
	// Dot is a symbol represented by any number of sequential '.' runes
	Dot
	// Semi is a ';' symbol
	Semi
	// Dollar is a '$' symbol
	Dollar
	// Comment is a token contained withing a `{` `}` characters or between `#` and `\n` characters
	// ... I think
	Comment

	// Quote is a symbole represented by text with in `"` runes
	Quote
	// Ident is a single word
	Ident
	// MoveNumber is a single positive non-zero whole number followed by a '.'
	MoveNumber
	// Number is a single positive number that may not be followed by a '.'
	Number

	// Nag represents a numeric algebraic glyph: https://en.wikipedia.org/wiki/Numeric_Annotation_Glyphs
	Nag
	// Result is a game result either *, 1/2-1/2, 1-0 or 0-1
	Result
)

// Token represents an atom of a pgn file
type Token struct {
	Tok      Tok
	Position scanner.Position
	Length   int
	Literal  string
}

// Scanner is a pgn scanner. For most applications is it recommended to use
// a pgn parser instead of the scanner
type Scanner struct {
	s scanner.Scanner
}

// Peek returns the next character in the input passing over whitespace, but does not move
// the next token pointer along.
func (ps *Scanner) Peek() rune {
	for {
		var char = ps.s.Peek()
		if !isWhitespace(char) {
			return char
		}
		ps.s.Next()
	}
}

// Next return the next token in the input passing over white space
func (ps *Scanner) Next() Token {
	for {
		tok := ps.next()
		if tok.Tok != Ws {
			return tok
		}
	}
}

func (ps *Scanner) next() Token {
	char := ps.s.Peek()

	if isWhitespace(char) {
		return ps.scanWhitespace()
	} else if isLetter(char) {
		return ps.scanIdent()
	} else if '"' == char {
		return ps.scanDoubleQuoted()
	} else if isNumber(char) {
		return ps.scanNumber()
	} else if '{' == char {
		return ps.scanComment()
	}

	ps.s.Next()

	switch char {
	case eof:
		return Token{
			Tok:      EOF,
			Position: ps.s.Pos(),
			Length:   1,
		}
	case '[':
		return Token{
			Tok:      LBrace,
			Position: ps.s.Pos(),
			Length:   1,
		}
	case ']':
		return Token{
			Tok:      RBrace,
			Position: ps.s.Pos(),
			Length:   1,
		}
	case '(':
		return Token{
			Tok:      LParen,
			Position: ps.s.Pos(),
			Length:   1,
		}
	case ')':
		return Token{
			Tok:      RParen,
			Position: ps.s.Pos(),
			Length:   1,
		}
	case '!', '?', '‼', '⁇', '⁉', '⁈', '□', '=', '∞', '±', '∓', '+', '-', '⨀', '⟳', '→', '↑', '⇆':
		str := string(char)
		length := 1

		for {
			if !isNag(ps.s.Peek()) {
				break
			}
			next := ps.s.Next()
			length++
			str += string(next)
		}

		return Token{
			Tok:      Nag,
			Position: ps.s.Pos(),
			Literal:  str,
			Length:   length,
		}
	case '*':
		return Token{
			Tok:      Result,
			Position: ps.s.Pos(),
			Length:   1,
			Literal:  "*",
		}
	}

	return Token{Tok: Illegal, Position: ps.s.Pos(), Length: 1}
}

// Init initializes a scanner with an io reader
func (ps *Scanner) Init(r io.Reader) {
	ps.s.Init(r)
}
