package optimizer

import "express_compiler/parser"

// strengthReduce replaces expensive operations with cheaper equivalents
func (o *Optimizer) strengthReduce(expr parser.Expr) parser.Expr {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *parser.BinaryExpr:
		// First reduce children
		e.Left = o.strengthReduce(e.Left)
		e.Right = o.strengthReduce(e.Right)

		// Then try to reduce this node
		return o.reduceStrength(e)

	case *parser.UnaryExpr:
		e.Right = o.strengthReduce(e.Right)
		return e

	case *parser.CallExpr:
		for i, arg := range e.Args {
			e.Args[i] = o.strengthReduce(arg)
		}
		return o.reduceCallStrength(e)

	default:
		return e
	}
}

func (o *Optimizer) reduceStrength(expr *parser.BinaryExpr) parser.Expr {
	switch expr.Operator {
	case "*":
		return o.reduceMultiplication(expr)
	case "/":
		return o.reduceDivision(expr)
	case "^":
		return o.reducePower(expr)
	default:
		return expr
	}
}

func (o *Optimizer) reduceMultiplication(expr *parser.BinaryExpr) parser.Expr {
	// x * 2 → x + x
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 2 {
		o.stats.StrengthReduced++
		return &parser.BinaryExpr{
			Left:     expr.Left,
			Operator: "+",
			Right:    expr.Left,
			Pos:      expr.Pos,
		}
	}

	// 2 * x → x + x
	if leftInt, ok := expr.Left.(*parser.IntLiteral); ok && leftInt.Value == 2 {
		o.stats.StrengthReduced++
		return &parser.BinaryExpr{
			Left:     expr.Right,
			Operator: "+",
			Right:    expr.Right,
			Pos:      expr.Pos,
		}
	}

	// x * 4 → (x + x) + (x + x)
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 4 {
		o.stats.StrengthReduced++
		xPlusX := &parser.BinaryExpr{
			Left:     expr.Left,
			Operator: "+",
			Right:    expr.Left,
			Pos:      expr.Pos,
		}
		return &parser.BinaryExpr{
			Left:     xPlusX,
			Operator: "+",
			Right:    xPlusX,
			Pos:      expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) reduceDivision(expr *parser.BinaryExpr) parser.Expr {
	// x / 2 could be reduced to x >> 1 for integers, but we don't have bit shifts
	// So we leave division as-is for now
	return expr
}

func (o *Optimizer) reducePower(expr *parser.BinaryExpr) parser.Expr {
	// x ^ 2 → x * x
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok {
		if rightInt.Value == 2 {
			o.stats.StrengthReduced++
			return &parser.BinaryExpr{
				Left:     expr.Left,
				Operator: "*",
				Right:    expr.Left,
				Pos:      expr.Pos,
			}
		}

		// x ^ 3 → x * x * x
		if rightInt.Value == 3 {
			o.stats.StrengthReduced++
			xTimesX := &parser.BinaryExpr{
				Left:     expr.Left,
				Operator: "*",
				Right:    expr.Left,
				Pos:      expr.Pos,
			}
			return &parser.BinaryExpr{
				Left:     xTimesX,
				Operator: "*",
				Right:    expr.Left,
				Pos:      expr.Pos,
			}
		}
	}

	return expr
}

func (o *Optimizer) reduceCallStrength(expr *parser.CallExpr) parser.Expr {
	// sqrt(x * x) → abs(x)  (not exactly true for negative x, but close)
	// For now, we don't do function strength reduction
	return expr
}
