package pgn

import (
	"bytes"
	"strings"
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumber(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isIdentChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '=' ||
		ch == '+' ||
		ch == '#' ||
		ch == '-'
}

func (ps *Scanner) scanWhitespace() Token {
	length := 0
	pos := ps.s.Pos()

	for {
		var char = ps.s.Peek()
		if !isWhitespace(char) {
			break
		} else {
			length++
			ps.s.Next()
		}
	}

	return Token{
		Tok:      Ws,
		Position: pos,
		Length:   length,
	}
}

func (ps *Scanner) scanIdent() Token {
	length := 0
	pos := ps.s.Pos()
	var buf bytes.Buffer

	for {
		var char = ps.s.Peek()
		if !isIdentChar(char) {
			break
		} else {
			length++
			ps.s.Next()
			_, _ = buf.WriteRune(char)
		}
	}

	return Token{
		Tok:      Ident,
		Position: pos,
		Length:   length,
		Literal:  buf.String(),
	}
}

func (ps *Scanner) scanNumber() Token {
	length := 1
	pos := ps.s.Pos()
	var str string

	firstChar := ps.s.Next()
	str += string(firstChar)

	secondChar := ps.s.Peek()
	switch secondChar {
	case '/':
		ch := ps.s.Next()
		str += string(ch)
		length++

		ch = ps.s.Next()
		str += string(ch)
		length++
		if ch != '2' {
			return Token{
				Tok:      Illegal,
				Length:   length,
				Position: pos,
				Literal:  str,
			}
		}

		ch = ps.s.Next()
		str += string(ch)
		length++
		if ch != '-' {
			return Token{
				Tok:      Illegal,
				Length:   length,
				Position: pos,
				Literal:  str,
			}
		}

		ch = ps.s.Next()
		str += string(ch)
		length++
		if ch != '1' {
			return Token{
				Tok:      Illegal,
				Length:   length,
				Position: pos,
				Literal:  str,
			}
		}

		ch = ps.s.Next()
		str += string(ch)
		length++
		if ch != '/' {
			return Token{
				Tok:      Illegal,
				Length:   length,
				Position: pos,
				Literal:  str,
			}
		}

		ch = ps.s.Next()
		str += string(ch)
		length++
		if ch != '2' {
			return Token{
				Tok:      Illegal,
				Length:   length,
				Position: pos,
				Literal:  str,
			}
		}

		return Token{
			Tok:      Result,
			Length:   length,
			Position: pos,
			Literal:  str,
		}
	case '-':
		ch := ps.s.Next()
		str += string(ch)
		length++
		if ch != '-' {
			return Token{
				Tok:      Illegal,
				Length:   length,
				Position: pos,
				Literal:  str,
			}
		}

		ch = ps.s.Next()
		str += string(ch)
		length++
		if ch != '0' && ch != '1' {
			return Token{
				Tok:      Illegal,
				Length:   length,
				Position: pos,
				Literal:  str,
			}
		}

		return Token{
			Tok:      Result,
			Length:   length,
			Position: pos,
			Literal:  str,
		}

	case '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		for {
			c := ps.s.Peek()
			if c == '.' {
				ps.s.Next()
				length++
				for {
					if ps.s.Peek() != '.' {
						return Token{
							Tok:      MoveNumber,
							Length:   length,
							Position: pos,
							Literal:  str,
						}
					}
					ps.s.Next()
					length++
				}
			}
			if c >= '0' && c <= '9' {
				next := ps.s.Next()

				str += string(next)
				length++
			} else {
				return Token{
					Tok:      Number,
					Length:   length,
					Position: pos,
					Literal:  str,
				}
			}
		}
	}

	return Token{}
}

func (ps *Scanner) scanDoubleQuoted() Token {
	ps.s.Next()
	length := 0
	pos := ps.s.Pos()
	var buf bytes.Buffer

	for {
		var char = ps.s.Peek()
		if char == '\\' {
			length++
			ps.s.Next()
			c := ps.s.Next()
			_, _ = buf.WriteRune(c)
		} else if char == '"' {
			ps.s.Next()
			break
		} else {
			length++
			ps.s.Next()
			_, _ = buf.WriteRune(char)
		}
	}

	return Token{
		Tok:      Quote,
		Position: pos,
		Length:   length,
		Literal:  buf.String(),
	}
}

func (ps *Scanner) scanComment() Token {
	ps.s.Next()
	length := 0
	pos := ps.s.Pos()
	var buf bytes.Buffer

	for {
		var char = ps.s.Peek()

		if char == '}' {
			ps.s.Next()
			break
		} else {
			length++
			ps.s.Next()
			_, _ = buf.WriteRune(char)
		}
	}

	return Token{
		Tok:      Comment,
		Position: pos,
		Length:   length,
		Literal:  strings.Trim(buf.String(), " "),
	}
}

func isNag(ch rune) bool {
	switch ch {
	case '!', '?', '‼', '⁇', '⁉', '⁈', '□', '=', '∞', '±', '∓', '+', '-', '⨀', '⟳', '→', '↑', '⇆':
		return true
	default:
		return false
	}
}
