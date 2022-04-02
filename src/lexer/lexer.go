package lexer

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"github.com/maxivhuber/monkeyi/token"
)

type Lexer struct {
	input        []byte
	position     int  // current position in input (byte offset to current rune)
	readPosition int  // current reading position in input (byte offset after current rune)
	rn           rune // current char under examination
}

func New(input string) (l *Lexer, err error) {

	if i := []byte(input); !utf8.Valid(i) {
		return nil, errors.New("invalid UTF-8 input")
	} else {
		l := &Lexer{input: i}

		l.readRune()
		return l, nil
	}

}

func (l *Lexer) readRune() {

	if l.readPosition >= len(l.input) {
		l.rn = 0
	} else {
		r, size := utf8.DecodeRune(l.input[l.readPosition:])

		l.position = l.readPosition
		l.readPosition += size
		l.rn = r
	}

}
func (l *Lexer) NextToken() token.Token {

	var tok token.Token

	l.eatWhitespace()

	switch l.rn {
	case '=':
		t, err := l.makeTwoRuneToken('=')
		if err == nil {
			tok = t
		} else {
			tok = newToken(token.ASSIGN, l.rn)
		}
	case '+':
		tok = newToken(token.PLUS, l.rn)
	case '-':
		tok = newToken(token.MINUS, l.rn)
	case '!':
		t, err := l.makeTwoRuneToken('!')
		if err == nil {
			tok = t
		} else {
			tok = newToken(token.BANG, l.rn)
		}

	case '/':
		tok = newToken(token.SLASH, l.rn)
	case '*':
		tok = newToken(token.ASTERISK, l.rn)
	case '<':
		tok = newToken(token.LT, l.rn)
	case '>':
		tok = newToken(token.GT, l.rn)
	case ';':
		tok = newToken(token.SEMICOLON, l.rn)
	case ',':
		tok = newToken(token.COMMA, l.rn)
	case '{':
		tok = newToken(token.LBRACE, l.rn)
	case '}':
		tok = newToken(token.RBRACE, l.rn)
	case '(':
		tok = newToken(token.LPAREN, l.rn)
	case ')':
		tok = newToken(token.RPAREN, l.rn)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		switch {
		case unicode.IsLetter(l.rn):
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		case unicode.IsDigit(l.rn):
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		default:
			tok = newToken(token.ILLEGAL, l.rn)
		}

	}

	l.readRune()
	return tok

}

func (l *Lexer) readIdentifier() string {

	ident := []rune{l.rn}

	for {
		l.readRune()
		if !unicode.IsLetter(l.rn) {
			break
		}
		ident = append(ident, l.rn)
	}
	return string(ident)

}

func (l *Lexer) readNumber() string {

	ident := []rune{l.rn}

	for {
		l.readRune()
		if !unicode.IsDigit(l.rn) {
			break
		}
		ident = append(ident, l.rn)
	}
	return string(ident)

}

func (l *Lexer) peekRune() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		rn, _ := utf8.DecodeRune(l.input[l.readPosition:])
		return rn
	}
}

func (l *Lexer) eatWhitespace() {

	for unicode.IsSpace(l.rn) {
		l.readRune()
	}
}

func (l *Lexer) makeTwoRuneToken(rn rune) (t token.Token, err error) {

	switch rn {
	case '=':
		if l.peekRune() == '=' {
			return l.newTwoRuneToken(token.EQ, rn), nil
		}
	case '!':
		if l.peekRune() == '=' {
			return l.newTwoRuneToken(token.NOT_EQ, rn), nil
		}
	}

	return token.Token{}, errors.New("One rune token")

}

func newToken(tokenType token.TokenType, rn rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(rn)}
}

func (l *Lexer) newTwoRuneToken(tokenType token.TokenType, rn rune) token.Token {
	l.readRune()
	literal := string(rn) + string(l.rn)
	return token.Token{Type: tokenType, Literal: literal}
}
