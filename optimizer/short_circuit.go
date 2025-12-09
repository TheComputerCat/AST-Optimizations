package optimizer

import "express_compiler/parser"

// shortCircuit optimizes boolean expressions using short-circuit logic
func (o *Optimizer) shortCircuit(expr parser.Expr) parser.Expr {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *parser.BinaryExpr:
		// First optimize children
		e.Left = o.shortCircuit(e.Left)
		e.Right = o.shortCircuit(e.Right)

		// Then try to short-circuit this node
		return o.shortCircuitBinary(e)

	case *parser.UnaryExpr:
		e.Right = o.shortCircuit(e.Right)
		return e

	case *parser.CallExpr:
		for i, arg := range e.Args {
			e.Args[i] = o.shortCircuit(arg)
		}
		return e

	default:
		return e
	}
}

func (o *Optimizer) shortCircuitBinary(expr *parser.BinaryExpr) parser.Expr {
	switch expr.Operator {
	case "&&":
		return o.shortCircuitAnd(expr)
	case "||":
		return o.shortCircuitOr(expr)
	default:
		return expr
	}
}

func (o *Optimizer) shortCircuitAnd(expr *parser.BinaryExpr) parser.Expr {
	// false && x → false
	if leftBool, ok := expr.Left.(*parser.BoolLiteral); ok {
		if !leftBool.Value {
			o.stats.ShortCircuited++
			return &parser.BoolLiteral{Value: false, Pos: expr.Pos}
		}
	}

	// x && false → false
	if rightBool, ok := expr.Right.(*parser.BoolLiteral); ok {
		if !rightBool.Value {
			o.stats.ShortCircuited++
			return &parser.BoolLiteral{Value: false, Pos: expr.Pos}
		}
	}

	// true && x → x
	if leftBool, ok := expr.Left.(*parser.BoolLiteral); ok {
		if leftBool.Value {
			o.stats.ShortCircuited++
			return expr.Right
		}
	}

	// x && true → x
	if rightBool, ok := expr.Right.(*parser.BoolLiteral); ok {
		if rightBool.Value {
			o.stats.ShortCircuited++
			return expr.Left
		}
	}

	return expr
}

func (o *Optimizer) shortCircuitOr(expr *parser.BinaryExpr) parser.Expr {
	// true || x → true
	if leftBool, ok := expr.Left.(*parser.BoolLiteral); ok {
		if leftBool.Value {
			o.stats.ShortCircuited++
			return &parser.BoolLiteral{Value: true, Pos: expr.Pos}
		}
	}

	// x || true → true
	if rightBool, ok := expr.Right.(*parser.BoolLiteral); ok {
		if rightBool.Value {
			o.stats.ShortCircuited++
			return &parser.BoolLiteral{Value: true, Pos: expr.Pos}
		}
	}

	// false || x → x
	if leftBool, ok := expr.Left.(*parser.BoolLiteral); ok {
		if !leftBool.Value {
			o.stats.ShortCircuited++
			return expr.Right
		}
	}

	// x || false → x
	if rightBool, ok := expr.Right.(*parser.BoolLiteral); ok {
		if !rightBool.Value {
			o.stats.ShortCircuited++
			return expr.Left
		}
	}

	return expr
}
