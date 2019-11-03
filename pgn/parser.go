package pgn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type parser struct {
	p Scanner
}

func newParser(p Scanner) parser {
	return parser{p}
}

type gameError struct {
	game Game
	err  error
}

// ParseConcurrent provides a concurrent structure to parse pgn. This
// function accpets an os.File instead of the more generic i.Reader interface.
func ParseConcurrent(filePath string) ([]Game, error) {
	var games []Game
	// My processor can handle about 25000 games per second
	// so I will break my file up in to sections of that size
	// and let a new go routine handle each section
	gamesPerRoutine := int64(100000)
	// There are roughly 800 runes per game
	charsPerGame := int64(800)
	// the problem will be broken up into sections that should take
	// about a second to process
	runesPerSection := gamesPerRoutine * charsPerGame

	// Open the file to get the file size
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Get the info for the file to see how long it is
	info, err := file.Stat()
	if err != nil {
		return games, err
	}

	// We need on extra chunk for the last section
	numChunks := int(info.Size()/runesPerSection) + 1

	// A channel to collect games on
	inChan := make(chan []Game)

	// processGameSection is a function that will process a number of games
	// up to a certain point in parser
	// The isStart flag indicates if we are starting from the beginning of a file because slightly different
	// logic applies in this case for getting the first token
	processGamesSection := func(filePath string, numRunes int64, start int64, outChan chan []Game) {
		var games []Game
		var scanner Scanner

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.Seek(start, os.SEEK_SET)
		if err != nil {
			fmt.Println(err)
			return
		}
		scanner.Init(file)
		p := parser{scanner}

		p.recover(start == 0)

		for {
			if p.p.Peek() == eof || p.p.s.Pos().Offset > int(numRunes) {
				break
			}

			game, err := p.parseGame()
			if err != nil {
				fmt.Println(err)
			} else {
				games = append(games, game)
			}
		}

		fmt.Println("Finished")
		outChan <- games
		return
	}

	for i := 0; i < numChunks; i++ {
		i64 := int64(i)
		start := i64 * runesPerSection

		go processGamesSection(filePath, runesPerSection, start, inChan)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// count keeps track of the number of routines that have finished
	count := 0
Loop:
	for {
		select {
		case newGames := <-inChan:
			games = append(games, newGames...)
			count++
			if count == numChunks {
				break Loop
			}
		case <-ctx.Done():
			break Loop
		}
	}

	return games, nil
}

// Parse is the recommended way of parsing a pgn file into a list of games.
func Parse(in io.Reader) ([]Game, error) {
	var s Scanner

	s.Init(in)

	p := parser{s}

	return p.parsePgn()
}

// Recovers moves ahead to the next game starting position.
// If a game has invalid format recover will skip ahead till the next
// game and continue parsing.
func (p *parser) recover(isStart bool) {
	if isStart {
		return
	}
	// We need to find two new lines folled by a [
	// and read all the way up to but not including the brace
	for {
		char1 := p.p.s.Next()
		if char1 != '\n' {
			continue
		}
		char2 := p.p.s.Peek()
		if char2 != '\n' {
			continue
		}
		for {
			char3 := p.p.s.Peek()
			if char3 != '\n' {
				break
			}
			p.p.s.Next()
		}
		if p.p.s.Peek() == '[' {
			break
		}
	}
}

// parsePgn takes a full pgn file and returns a list of games within that file
// Note: this function is not concurrent. For a concurrent version use the pgn.Parse function
func (p *parser) parsePgn() ([]Game, error) {
	var errorList []error
	var games []Game

	for {
		if p.p.Peek() == eof {
			break
		}

		game, err := p.parseGame()
		if err != nil {
			errorList = append(errorList, err)
			p.recover(false)
		} else {
			games = append(games, game)
		}
	}

	return games, consolicateErrors("parsing pgn", errorList)
}

func (p *parser) parseGame() (Game, error) {
	var game Game

	tags, err := p.parseTags()
	if err != nil {
		return game, err
	}
	game.Tags = tags

	moves, err := p.parseMoves()
	if err != nil {
		return game, err
	}
	game.Moves = moves

	return game, nil
}

func (p *parser) parseTags() (map[string]string, error) {
	tags := make(map[string]string)

	for {
		lbrace := p.p.Next()
		if lbrace.Tok != LBrace {
			return tags, invalidToken("[", lbrace)
		}

		key := p.p.Next()
		if key.Tok != Ident {
			return tags, invalidToken("IDENT", key)
		}

		value := p.p.Next()
		if value.Tok != Quote {
			return tags, invalidToken("QUOTE", value)
		}

		rbrace := p.p.Next()
		if rbrace.Tok != RBrace {
			return tags, invalidToken("]", rbrace)
		}

		tags[key.Literal] = value.Literal

		if p.p.Peek() != '[' {
			break
		}
	}

	return tags, nil
}

// parseMoves must contain at least on move and the first move must contain a move number
// apart from that all subsequent moves must only contain a move string and may optional contain a
// move number, nag, annotation, or alternate moves
// TODO: moves with diagrams
func (p *parser) parseMoves() ([]Move, error) {
	var moves []Move
	tok := p.p.Next()

	// Some games have no moves
	if tok.Tok == Result {
		return moves, nil
	}

Loop:
	for {
		// Init move
		var move Move

		// Check for move number
		if tok.Tok == MoveNumber {
			i, err := strconv.Atoi(tok.Literal)
			if err != nil {
				return moves, fmt.Errorf("Could not convert string to integer in move number. \n%v", tok.Position)
			}

			move.Number = int32(i)
			tok = p.p.Next()
		}

		// There must be a move string next
		if tok.Tok != Ident {
			return moves, invalidToken("Ident", tok)
		}
		move.Move = tok.Literal

		tok = p.p.Next()
		// Check for nag
		if tok.Tok == Nag {
			move.Nag = tok.Literal
			tok = p.p.Next()
		}

		// Check for comment
		if tok.Tok == Comment {
			move.Annotation = strings.Trim(tok.Literal, " ")
			tok = p.p.Next()
		}

		// Check for alternative
		if tok.Tok == LParen {
			return moves, errors.New("Alternatives is unimplemented")
		}

		moves = append(moves, move)
		// Check for result
		if tok.Tok == Result {
			break Loop
		}

	}

	return moves, nil
}

// parseMoveStr returns a move string if the next token on scanner is a valid move and an error that is not the case
// TODO: Some one more sober (maybe future Banner) check this is sound and complete... I don't think the tests cover all cases.
func (p *parser) parseMoveStr() (string, error) {
	tok := p.p.Next()
	if tok.Tok != Ident {
		return "", invalidToken("move", tok)
	}

	if tok.Literal == "O-O-O" || tok.Literal == "O-O" {
		return tok.Literal, nil
	}

	firstChar := tok.Literal[0]
	if firstChar != 'N' && firstChar != 'B' && firstChar != 'R' && firstChar != 'Q' && firstChar != 'K' && (firstChar < 'a' || firstChar > 'h') {
		return "", invalidToken("move", tok)
	}

	for _, ch := range tok.Literal {
		passed := false
		for _, a := range "abcdefghNBRQK12345678=+x-#" {
			if ch == a {
				passed = true
				break
			}
		}

		if !passed {
			return "", invalidToken("move", tok)
		}
	}

	return tok.Literal, nil
}

func invalidToken(want string, got Token) error {
	return fmt.Errorf("Invalid token on: %v\n\texpecting: \"%v\" but got: \"%s\"", got.Position, want, got.Literal)
}

func consolicateErrors(msg string, errorList []error) error {
	if len(errorList) == 0 {
		return nil
	}

	var errs string
	for i, err := range errorList[0:10] {
		errs += fmt.Sprintf("\t%v. %v\n\n", i, err)
	}
	if len(errorList) > 10 {
		errs += fmt.Sprintf("And %v other errors", len(errorList)-10)
	}
	return fmt.Errorf("%v errors occured when %s\n\n%s", len(errorList), msg, errs)
}
