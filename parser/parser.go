package parser

import (
	"express_compiler/lexer"
	"fmt"
	"strconv"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	LOGIC_OR    // ||
	LOGIC_AND   // &&
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SUM         // + -
	PRODUCT     // * / %
	POWER       // ^
	PREFIX      // -x !x
	REGEX       // ~
	CALL        // func(x)
)

var precedences = map[lexer.Type]int{
	lexer.OR:       LOGIC_OR,
	lexer.AND:      LOGIC_AND,
	lexer.EQ:       EQUALS,
	lexer.NEQ:      EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.LTE:      LESSGREATER,
	lexer.GTE:      LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTARISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
	lexer.MODULO:   PRODUCT,
	lexer.POWER:    POWER,
	lexer.REGEX:    REGEX,
	lexer.LPAREN:   CALL,
}

type Parser struct {
	l         lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(l lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.Type) {
	msg := fmt.Sprintf("line %d:%d - expected next token to be %s, got %s instead",
		p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t lexer.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseExpression is the main entry point
func (p *Parser) ParseExpression() Expr {
	return p.parseExpression(LOWEST)
}

func (p *Parser) parseExpression(precedence int) Expr {
	// Parse prefix expression
	left := p.parsePrefixExpression()
	if left == nil {
		return nil
	}

	// Parse infix expressions with precedence climbing
	for !p.peekTokenIs(lexer.EOF) && precedence < p.peekPrecedence() {
		p.nextToken()
		left = p.parseInfixExpression(left)
	}

	return left
}

func (p *Parser) parsePrefixExpression() Expr {
	switch p.curToken.Type {
	case lexer.INT:
		return p.parseIntLiteral()
	case lexer.FLOAT:
		return p.parseFloatLiteral()
	case lexer.STRING:
		return p.parseStringLiteral()
	case lexer.TRUE, lexer.FALSE:
		return p.parseBoolLiteral()
	case lexer.IDENT:
		// Could be identifier or function call
		if p.peekTokenIs(lexer.LPAREN) {
			return p.parseCallExpression()
		}
		return p.parseIdentifier()
	case lexer.NOT, lexer.MINUS:
		return p.parseUnaryExpression()
	case lexer.LPAREN:
		return p.parseGroupedExpression()
	default:
		msg := fmt.Sprintf("line %d:%d - no prefix parse function for %s",
			p.curToken.Line, p.curToken.Column, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
}

func (p *Parser) parseInfixExpression(left Expr) Expr {
	if p.curTokenIs(lexer.LPAREN) && left != nil {
		// This is a call expression
		if ident, ok := left.(*Identifier); ok {
			return p.parseCallExpressionWithIdent(ident)
		}
	}

	expr := &BinaryExpr{
		Left:     left,
		Operator: p.curToken.Literal,
		Pos:      Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}

	precedence := p.curPrecedence()

	// Right-associative for power operator
	if p.curToken.Type == lexer.POWER {
		precedence-- // Use lower precedence for right side
	}

	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseIntLiteral() Expr {
	lit := &IntLiteral{
		Pos: Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("line %d:%d - could not parse %q as integer",
			p.curToken.Line, p.curToken.Column, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expr {
	lit := &FloatLiteral{
		Pos: Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("line %d:%d - could not parse %q as float",
			p.curToken.Line, p.curToken.Column, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expr {
	return &StringLiteral{
		Value: p.curToken.Literal,
		Pos:   Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}
}

func (p *Parser) parseBoolLiteral() Expr {
	return &BoolLiteral{
		Value: p.curTokenIs(lexer.TRUE),
		Pos:   Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}
}

func (p *Parser) parseIdentifier() Expr {
	return &Identifier{
		Name: p.curToken.Literal,
		Pos:  Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}
}

func (p *Parser) parseUnaryExpression() Expr {
	expr := &UnaryExpr{
		Operator: p.curToken.Literal,
		Pos:      Position{Line: p.curToken.Line, Column: p.curToken.Column},
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseGroupedExpression() Expr {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseCallExpression() Expr {
	funcName := p.curToken.Literal
	pos := Position{Line: p.curToken.Line, Column: p.curToken.Column}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	args := p.parseCallArguments()

	return &CallExpr{
		Function: funcName,
		Args:     args,
		Pos:      pos,
	}
}

func (p *Parser) parseCallExpressionWithIdent(ident *Identifier) Expr {
	args := p.parseCallArguments()

	return &CallExpr{
		Function: ident.Name,
		Args:     args,
		Pos:      ident.Pos,
	}
}

func (p *Parser) parseCallArguments() []Expr {
	args := []Expr{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return args
}
