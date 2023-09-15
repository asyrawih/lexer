package main

import (
	"fmt"
	"log"
	"strconv"
)

type Node interface {
	Pos() Position
	String() string
}

type Expression interface {
	Node
	exprNode()
}

type BinaryExpression struct {
	Left     Expression
	Op       Token
	Right    Expression
	Position Position
}

func (be *BinaryExpression) Pos() Position {
	return be.Position
}

func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), tokens[be.Op], be.Right.String())
}

func (be *BinaryExpression) exprNode() {}

type IntegerLiteral struct {
	Value    int
	Position Position
}

// exprNode implements Expression.
func (*IntegerLiteral) exprNode() {
	panic("unimplemented")
}

func (il *IntegerLiteral) Pos() Position {
	return il.Position
}

func (il *IntegerLiteral) String() string {
	return strconv.Itoa(il.Value)
}

func parseExpression(l *Lexer) Expression {
	return parseAddSubExpr(l)
}

func parseAddSubExpr(l *Lexer) Expression {
	left := parseMulDivExpr(l)

	for {
		tok, _ := l.Lex()
		if tok != ADD && tok != SUB {
			l.backup()
			return left
		}

		right := parseMulDivExpr(l)
		left = &BinaryExpression{Left: left, Op: tok, Right: right, Position: left.Pos()}
	}
}

func parseMulDivExpr(l *Lexer) Expression {
	left := parsePrimaryExpr(l)

	for {
		tok, _ := l.Lex()
		if tok != MUL && tok != DIV {
			l.backup()
			return left
		}

		right := parsePrimaryExpr(l)
		left = &BinaryExpression{Left: left, Op: tok, Right: right, Position: left.Pos()}
	}
}

func parsePrimaryExpr(l *Lexer) Expression {
	tok, lit := l.Lex()

	if tok == INT {
		value, _ := strconv.Atoi(lit)
		return &IntegerLiteral{Value: value, Position: Position{line: l.pos.line, column: l.pos.column}}
	}

	log.Fatalf("Unexpected token: %s", tokens[tok])
	return nil // unreachable
}

func evaluateExpression(expr Expression) (int, error) {
	switch e := expr.(type) {
	case *BinaryExpression:
		left, err := evaluateExpression(e.Left)
		if err != nil {
			return 0, err
		}

		right, err := evaluateExpression(e.Right)
		if err != nil {
			return 0, err
		}

		switch e.Op {
		case ADD:
			return left + right, nil
		case SUB:
			return left - right, nil
		case MUL:
			return left * right, nil
		case DIV:
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("unknown operator")
		}

	case *IntegerLiteral:
		return e.Value, nil

	default:
		return 0, fmt.Errorf("unknown expression type")
	}
}
