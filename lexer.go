package main

import (
	"bufio"
	"io"
	"log"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	ADD:     "+",
	SUB:     "-",
	MUL:     "*",
	DIV:     "/",
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func (l *Lexer) Lex() (Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return EOF, ""
			}
			log.Fatal(err)
		}
		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case '+':
			return ADD, "+"
		case '-':
			return SUB, "-"
		case '*':
			return MUL, "*"
		case '/':
			return DIV, "/"
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				l.backup()
				lit := l.lexInt()
				return INT, lit
			} else {
				return ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		log.Fatal(err)
	}
	l.pos.column--
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
			log.Fatal(err)
		}
		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func NewLexer(reader *bufio.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}
