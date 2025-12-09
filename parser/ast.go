package parser

import "fmt"

type Node interface {
	String() string
}

type Expr interface {
	Node
	exprNode()
}

type Position struct {
	Line   int
	Column int
}

type IntLiteral struct {
	Value int64
	Pos   Position
}

func (i *IntLiteral) exprNode() {}
func (i *IntLiteral) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type FloatLiteral struct {
	Value float64
	Pos   Position
}

func (f *FloatLiteral) exprNode() {}
func (f *FloatLiteral) String() string {
	return fmt.Sprintf("%f", f.Value)
}

type StringLiteral struct {
	Value string
	Pos   Position
}

func (s *StringLiteral) exprNode() {}
func (s *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

type BoolLiteral struct {
	Value bool
	Pos   Position
}

func (b *BoolLiteral) exprNode() {}
func (b *BoolLiteral) String() string {
	return fmt.Sprintf("%t", b.Value)
}

type Identifier struct {
	Name string
	Pos  Position
}

func (i *Identifier) exprNode() {}
func (i *Identifier) String() string {
	return i.Name
}

type BinaryExpr struct {
	Left     Expr
	Operator string
	Right    Expr
	Pos      Position
}

func (b *BinaryExpr) exprNode() {}
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator, b.Left.String(), b.Right.String())
}

type UnaryExpr struct {
	Operator string
	Right    Expr
	Pos      Position
}

func (u *UnaryExpr) exprNode() {}
func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator, u.Right.String())
}

type CallExpr struct {
	Function string
	Args     []Expr
	Pos      Position
}

func (c *CallExpr) exprNode() {}
func (c *CallExpr) String() string {
	args := ""
	for i, arg := range c.Args {
		if i > 0 {
			args += ", "
		}
		args += arg.String()
	}
	return fmt.Sprintf("%s(%s)", c.Function, args)
}
