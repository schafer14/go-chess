package pgn

import (
	"strings"
	"testing"
)

func TestNext(t *testing.T) {

	for _, tester := range testers {
		var s Scanner

		s.Init(strings.NewReader(tester.phrase))

		for i, expected := range tester.tokens {
			got := s.Next()

			if got.Tok != expected.Tok {
				t.Errorf("Error in test %s token %d: expected tok %v but got %v", tester.name, i, expected.Tok, got.Tok)
			}

			if got.Literal != expected.Literal {
				t.Errorf("Error in test %s token %d: expected literal '%s' but got '%s'", tester.name, i, expected.Literal, got.Literal)
			}

			if expected.Length != 0 && expected.Length != got.Length {
				t.Errorf("Error in test %s token %d: expected length %v but got %v", tester.name, i, expected.Length, got.Length)
			}
		}
	}
}

type scannerTest struct {
	phrase string
	tokens []Token
	name   string
}

var testers = []scannerTest{
	scannerTest{
		name:   "Lichess example comment with multi nag",
		phrase: "27... Bd5+?! { [%eval #-2] }",
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Length:  5,
				Literal: "27",
			},
			Token{
				Tok:     Ident,
				Literal: "Bd5+",
			},
			Token{
				Tok:     Nag,
				Literal: "?!",
			},
			Token{
				Tok:     Comment,
				Literal: "[%eval #-2]",
			},
		},
	},
	scannerTest{
		name:   "Tag Happy Case",
		phrase: `[White "Fabiano Caruana"]`,
		tokens: []Token{
			Token{
				Tok: LBrace,
			},
			Token{
				Tok:     Ident,
				Literal: "White",
			},
			Token{
				Tok:     Quote,
				Literal: "Fabiano Caruana",
			},
			Token{
				Tok: RBrace,
			},
			Token{
				Tok: EOF,
			},
		},
	},
	scannerTest{
		name:   `Tag with \\`,
		phrase: `[White "Fabiano \\Caruana"]`,
		tokens: []Token{
			Token{
				Tok: LBrace,
			},
			Token{
				Tok:     Ident,
				Literal: "White",
			},
			Token{
				Tok:     Quote,
				Literal: `Fabiano \Caruana`,
			},
			Token{
				Tok: RBrace,
			},
			Token{
				Tok: EOF,
			},
		},
	},
	scannerTest{
		name:   `Tag with \"`,
		phrase: `[White "Fabiano \"Caruana"]`,
		tokens: []Token{
			Token{
				Tok: LBrace,
			},
			Token{
				Tok:     Ident,
				Literal: "White",
			},
			Token{
				Tok:     Quote,
				Literal: `Fabiano "Caruana`,
			},
			Token{
				Tok: RBrace,
			},
			Token{
				Tok: EOF,
			},
		},
	},
	scannerTest{
		name:   `Happy Moves`,
		phrase: `1. e4 c5 2. Nf6`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok:     Ident,
				Literal: `c5`,
			},
			Token{
				Tok:     MoveNumber,
				Literal: "2",
			},
			Token{
				Tok:     Ident,
				Literal: "Nf6",
			},
		},
	},
	scannerTest{
		name:   `Comments`,
		phrase: `1. e4 {Incredible}`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok:     Comment,
				Literal: `Incredible`,
			},
		},
	},
	scannerTest{
		name:   `Dot dot dot notation`,
		phrase: `1... c5`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
				Length:  4,
			},
			Token{
				Tok:     Ident,
				Literal: `c5`,
			},
		},
	},
	scannerTest{
		name:   `Alternatives`,
		phrase: `1. e4 (1. d4)`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},

			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok: LParen,
			},
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: "d4",
			},
			Token{
				Tok: RParen,
			},
		},
	},
	scannerTest{
		name:   `White wins`,
		phrase: `1. e4 1-0`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok:     Result,
				Literal: "1-0",
				Length:  3,
			},
		},
	},
	scannerTest{
		name:   `Black wins`,
		phrase: `1. e4 0-1`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok:     Result,
				Literal: "0-1",
				Length:  3,
			},
		},
	},
	scannerTest{
		name:   `Draw`,
		phrase: `1. e4 1/2-1/2`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok:     Result,
				Literal: "1/2-1/2",
				Length:  7,
			},
		},
	},
	scannerTest{
		name:   `Result in progress`,
		phrase: `1. e4 *`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok:     Result,
				Literal: "*",
				Length:  1,
			},
		},
	},
	scannerTest{
		name:   `Nag`,
		phrase: `1. e4 !`,
		tokens: []Token{
			Token{
				Tok:     MoveNumber,
				Literal: "1",
			},
			Token{
				Tok:     Ident,
				Literal: `e4`,
			},
			Token{
				Tok:     Nag,
				Literal: "!",
				Length:  1,
			},
		},
	},
}
