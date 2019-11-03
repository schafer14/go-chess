package pgn

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	game1 = `[Event "Rated Classical game"]
[Site "https://lichess.org/j1dkb5dw"]
[White "BFG9k"]
[Black "mamalak"]
[Result "1-0"]
[UTCDate "2012.12.31"]
[UTCTime "23:01:03"]
[WhiteElo "1639"]
[BlackElo "1403"]
[WhiteRatingDiff "+5"]
[BlackRatingDiff "-8"]
[ECO "C00"]
[Opening "French Defense: Normal Variation"]
[TimeControl "600+8"]
[Termination "Normal"]

1. e4 e6 2. d4 b6 3. a3 Bb7 4. Nc3 Nh6 5. Bxh6 gxh6 6. Be2 Qg5 7. Bg4 h5 8. Nf3 Qg6 9. Nh4 Qg5 10. Bxh5 Qxh4 11. Qf3 Kd8 12. Qxf7 Nc6 13. Qe8# 1-0`
)

func Test_parser_parseGame(t *testing.T) {
	type fields struct {
		p Scanner
	}
	tests := []struct {
		name    string
		phrase  string
		want    Game
		wantErr bool
	}{
		{
			name:   "Parse game",
			phrase: game1,
			want: Game{Tags: map[string]string{
				"Event":           "Rated Classical game",
				"Site":            "https://lichess.org/j1dkb5dw",
				"White":           "BFG9k",
				"Black":           "mamalak",
				"Result":          "1-0",
				"UTCDate":         "2012.12.31",
				"UTCTime":         "23:01:03",
				"WhiteElo":        "1639",
				"BlackElo":        "1403",
				"WhiteRatingDiff": "+5",
				"BlackRatingDiff": "-8",
				"ECO":             "C00",
				"Opening":         "French Defense: Normal Variation",
				"TimeControl":     "600+8",
				"Termination":     "Normal",
			}, Moves: []Move{
				Move{
					Number: 1,
					Move:   "e4",
				},
				Move{
					Move: "e6",
				},
				Move{
					Number: 2,
					Move:   "d4",
				},
				Move{
					Move: "b6",
				},
				Move{
					Number: 3,
					Move:   "a3",
				},
				Move{
					Move: "Bb7",
				},
				Move{
					Number: 4,
					Move:   "Nc3",
				},
				Move{
					Move: "Nh6",
				},
				Move{
					Number: 5,
					Move:   "Bxh6",
				},
				Move{
					Move: "gxh6",
				},
				Move{
					Number: 6,
					Move:   "Be2",
				},
				Move{
					Move: "Qg5",
				},
				Move{
					Number: 7,
					Move:   "Bg4",
				},
				Move{
					Move: "h5",
				},
				Move{
					Number: 8,
					Move:   "Nf3",
				},
				Move{
					Move: "Qg6",
				},
				Move{
					Number: 9,
					Move:   "Nh4",
				},
				Move{
					Move: "Qg5",
				},
				Move{
					Number: 10,
					Move:   "Bxh5",
				},
				Move{
					Move: "Qxh4",
				},
				Move{
					Number: 11,
					Move:   "Qf3",
				},
				Move{
					Move: "Kd8",
				},
				Move{
					Number: 12,
					Move:   "Qxf7",
				},
				Move{
					Move: "Nc6",
				},
				Move{
					Number: 13,
					Move:   "Qe8#",
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Scanner

			s.Init(strings.NewReader(tt.phrase))

			p := parser{s}

			game, err := p.parseGame()
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.parseGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(game.Tags, tt.want.Tags) {
				t.Errorf("parser.parseGame().Tags = %v, want %v", game.Tags, tt.want.Tags)
			}
			if !reflect.DeepEqual(game.Moves, tt.want.Moves) {
				t.Errorf("parser.parseGame().Moves = %v, want %v", game.Moves, tt.want.Moves)
			}
		})
	}
}

func Test_parser_parsePgn(t *testing.T) {
	var s Scanner
	file, err := os.Open("./example.pgn")
	if err != nil {
		t.Errorf("Could not open example file.\n %v", err)
	}
	defer file.Close()

	s.Init(file)

	p := parser{s}

	games, err := p.parsePgn()
	if err != nil {
		t.Errorf("Could not parse pgn example file\n. %v", err)
	}

	if len(games) != 121332 {
		t.Errorf("Returned wrong number of games from example pgn.\nExpecting 121,332 but got %v", len(games))
	}
}

func Test_Parse(t *testing.T) {
	file, err := os.Open("./example.pgn")
	if err != nil {
		t.Errorf("Could not open example file.\n %v", err)
	}
	defer file.Close()

	games, err := Parse(file)
	if err != nil {
		t.Errorf("Could not parse pgn example file\n. %v", err)
	}

	if len(games) != 121332 {
		t.Errorf("Returned wrong number of games from example pgn.\nExpecting 121,332 but got %v", len(games))
	}
}

func Test_parser_parseTags(t *testing.T) {
	tests := []struct {
		name    string
		phrase  string
		want    map[string]string
		wantErr bool
	}{
		{
			name:   "Single Tag",
			phrase: `[White "Fabiano Caruana"]`,
			want: map[string]string{
				"White": "Fabiano Caruana",
			},
			wantErr: false,
		},
		{
			name:   "Multiple Tags",
			phrase: `[White "Fabiano Caruana"] [Black "Hikaru Nakamura"]`,
			want: map[string]string{
				"White": "Fabiano Caruana",
				"Black": "Hikaru Nakamura",
			},
			wantErr: false,
		},
		{
			name:   "Tags with awkward spacing",
			phrase: "\n [ White  \"Fabiano Caruana\" ]   \n[\tBlack \t\n\"Hikaru Nakamura\"  \n] ",
			want: map[string]string{
				"White": "Fabiano Caruana",
				"Black": "Hikaru Nakamura",
			},
			wantErr: false,
		},
		{
			name:    "Fails on move",
			phrase:  `1. e4`,
			want:    map[string]string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Scanner

			s.Init(strings.NewReader(tt.phrase))

			p := parser{s}

			got, err := p.parseTags()
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.parseTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_parseMoves(t *testing.T) {
	tests := []struct {
		name    string
		phrase  string
		want    []Move
		wantErr bool
	}{
		{
			name:   "Single numbered move",
			phrase: `1. e4 1/2-1/2`,
			want: []Move{
				Move{
					Number: 1,
					Move:   "e4",
				},
			},
			wantErr: false,
		},
		{
			name:   "Single numbered move for black",
			phrase: `1... c5 *`,
			want: []Move{
				Move{
					Number: 1,
					Move:   "c5",
				},
			},
			wantErr: false,
		},
		{
			name:   "Multiple moves",
			phrase: `1. e4 c5 2. Nf3 1-0`,
			want: []Move{
				Move{
					Number: 1,
					Move:   "e4",
				},
				Move{
					Move: "c5",
				},
				Move{
					Number: 2,
					Move:   "Nf3",
				},
			},
			wantErr: false,
		},
		{
			name:   "Move without result",
			phrase: `1. e4`,
			want: []Move{
				Move{
					Number: 1,
					Move:   "e4",
				},
			},
			wantErr: true,
		},
		{
			name:   "Move with comment",
			phrase: `1. e4 { This is e4 }`,
			want: []Move{
				Move{
					Number:     1,
					Move:       "e4",
					Annotation: "This is e4",
				},
			},
			wantErr: true,
		},
		{
			name:   "Move with nag",
			phrase: `1. e4 !`,
			want: []Move{
				Move{
					Number: 1,
					Move:   "e4",
					Nag:    "!",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Scanner

			s.Init(strings.NewReader(tt.phrase))

			p := parser{s}

			got, err := p.parseMoves()
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.parseMoves() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseMoves() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_parseMove(t *testing.T) {
	tests := []struct {
		name    string
		phrase  string
		want    string
		wantErr bool
	}{
		{
			name:    "Simple pawn move",
			phrase:  `e4`,
			want:    "e4",
			wantErr: false,
		},
		{
			name:    "Simple piece move",
			phrase:  `Nf6`,
			want:    "Nf6",
			wantErr: false,
		},
		{
			name:    "Invalid piece move",
			phrase:  `nf6`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid dest move",
			phrase:  `F6`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid dest move with piece",
			phrase:  `NF6`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Castle kingside",
			phrase:  `O-O`,
			want:    "O-O",
			wantErr: false,
		},
		{
			name:    "Castle queenside",
			phrase:  `O-O-O`,
			want:    "O-O-O",
			wantErr: false,
		},
		{
			name:    "Invalid spacing",
			phrase:  `O- O -O`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Checking",
			phrase:  `Qf7+`,
			want:    "Qf7+",
			wantErr: false,
		},
		{
			name:    "Check mate",
			phrase:  `Qf7#`,
			want:    "Qf7#",
			wantErr: false,
		},
		{
			name:    "Weird symbol",
			phrase:  `Qf7L`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Capture",
			phrase:  `Nxe4`,
			want:    "Nxe4",
			wantErr: false,
		},
		{
			name:    "Promotion",
			phrase:  `d8=Q`,
			want:    "d8=Q",
			wantErr: false,
		},
		// TODO: check this actually the way lichess works
		{
			name:    "Capture promotion and check",
			phrase:  `dxc8=Q+`,
			want:    "dxc8=Q+",
			wantErr: false,
		},
		{
			name:    "Promoting to weird piece",
			phrase:  `c8=L`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Pawn capture",
			phrase:  `cxd4`,
			want:    "cxd4",
			wantErr: false,
		},
		{
			name:    "Spcifier stuff",
			phrase:  `R2d2`,
			want:    "R2d2",
			wantErr: false,
		},
		{
			name:    "specifier capture",
			phrase:  `Ncxd2`,
			want:    "Ncxd2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Scanner

			s.Init(strings.NewReader(tt.phrase))

			p := parser{s}

			got, err := p.parseMoveStr()
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.parseMoveStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseMoveStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_recover(t *testing.T) {
	data := `[WhiteElo "1499"]
[BlackElo "1733"]
[WhiteRatingDiff "+16"]
[BlackRatingDiff "-83"]
[ECO "A00"]
[Opening "Polish Opening"]
[TimeControl "300+2"]
[Termination "Time forfeit"]

1. b4 Nf6 2. e3 e5 3. Bb2 d6 4. Nf3 h6 5. Be2 Be7 6. Nc3 O-O 7. O-O a6 8. Ne1 Be6 9. Bf3 c6 10. d4 exd4 11. exd4 d5 12. Nd3 Qc7 13. Bc1 Nbd7 14. g3 Bd6 15. Nf4 Bxb4 16. Nxe6 fxe6 17. Bf4 Bd6 18. Bxh6 gxh6 19. Qd2 Kg7 20. Nd1 e5 21. c3 exd4 22. cxd4 Rae8 23. Ne3 Ne4 24. Bxe4 Rxe4 25. Nf5+ Rxf5 26. f3 Re7 27. Rf2 Ref7 28. Qd1 Kh8 29. Qd2 Bf8 30. Qe1 Bg7 31. Qe8+ Kh7 32. g4 Rf4 33. Re1 Nb6 34. Re6 Rf4f6 35. Re3 Nc4 36. Ree2 Qf4 37. g5 Qxg5+ 38. Kf1 Qf4 39. Qc8 1-0

[Event "Rated Classical game"]
[Site "https://lichess.org/3g18fmpk"]
[White "UKwildcats"]
[Black "palko09"]
[Result "1-0"]
[UTCDate "2013.01.02"]
[UTCTime "01:46:33"]
[WhiteElo "1571"]
[BlackElo "1524"]
[WhiteRatingDiff "+9"]
[BlackRatingDiff "-21"]
[ECO "B32"]
[Opening "Sicilian Defense: Franco-Sicilian Variation"]
[TimeControl "480+10"]
[Termination "Normal"]

1. e4 c5 2. d4 e6 3. Nf3 Nc6 4. d5 exd5 5. exd5 Na5 6. Qe2+ Be7 7. d6 Nc6 8. Nc3 Nf6 9. Bg5 O-O 10. dxe7 Nxe7 11. Bxf6 gxf6 12. O-O-O d5 13. Qd2 d4 14. Qh6 Qa5 15. Bd3 f5 16. Qg5+ Ng6 17. Nxd4 cxd4 18. Ne4 fxe4 19. Qxa5 exd3 20. Rxd3 Rd8 21. Rhd1 Bf5 1-0
`
	var s Scanner
	s.Init(strings.NewReader(data))
	p := parser{s}

	p.recover(false)
	game, err := p.parseGame()
	if err != nil {
		t.Errorf("Unexpected error when recovering\n%v", err)
	}
	if game.Moves[0].Move != "e4" {
		t.Errorf("Error recovering expected 'e4' but got '%v'", game.Moves[0].Move)
	}
}
