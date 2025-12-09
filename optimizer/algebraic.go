package optimizer

import "express_compiler/parser"

// algebraicSimplify applies mathematical identities to simplify expressions
func (o *Optimizer) algebraicSimplify(expr parser.Expr) parser.Expr {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *parser.BinaryExpr:
		// First simplify children (bottom-up)
		e.Left = o.algebraicSimplify(e.Left)
		e.Right = o.algebraicSimplify(e.Right)

		// Then try to simplify this node
		return o.simplifyBinaryExpr(e)

	case *parser.UnaryExpr:
		e.Right = o.algebraicSimplify(e.Right)
		return o.simplifyUnaryExpr(e)

	case *parser.CallExpr:
		for i, arg := range e.Args {
			e.Args[i] = o.algebraicSimplify(arg)
		}
		return e

	default:
		return e
	}
}

func (o *Optimizer) simplifyBinaryExpr(expr *parser.BinaryExpr) parser.Expr {
	switch expr.Operator {
	case "+":
		return o.simplifyAddition(expr)
	case "-":
		return o.simplifySubtraction(expr)
	case "*":
		return o.simplifyMultiplication(expr)
	case "/":
		return o.simplifyDivision(expr)
	case "^":
		return o.simplifyPower(expr)
	default:
		return expr
	}
}

func (o *Optimizer) simplifyAddition(expr *parser.BinaryExpr) parser.Expr {
	// x + 0 → x
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 0 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	// 0 + x → x
	if leftInt, ok := expr.Left.(*parser.IntLiteral); ok && leftInt.Value == 0 {
		o.stats.AlgebraicSimplified++
		return expr.Right
	}

	// Check for float zero as well
	if rightFloat, ok := expr.Right.(*parser.FloatLiteral); ok && rightFloat.Value == 0.0 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	if leftFloat, ok := expr.Left.(*parser.FloatLiteral); ok && leftFloat.Value == 0.0 {
		o.stats.AlgebraicSimplified++
		return expr.Right
	}

	return expr
}

func (o *Optimizer) simplifySubtraction(expr *parser.BinaryExpr) parser.Expr {
	// x - 0 → x
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 0 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	if rightFloat, ok := expr.Right.(*parser.FloatLiteral); ok && rightFloat.Value == 0.0 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	// x - x → 0 (if both are identical)
	if exprEqual(expr.Left, expr.Right) {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 0, Pos: expr.Pos}
	}

	return expr
}

func (o *Optimizer) simplifyMultiplication(expr *parser.BinaryExpr) parser.Expr {
	// x * 0 → 0
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 0 {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 0, Pos: expr.Pos}
	}

	// 0 * x → 0
	if leftInt, ok := expr.Left.(*parser.IntLiteral); ok && leftInt.Value == 0 {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 0, Pos: expr.Pos}
	}

	// x * 1 → x
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 1 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	// 1 * x → x
	if leftInt, ok := expr.Left.(*parser.IntLiteral); ok && leftInt.Value == 1 {
		o.stats.AlgebraicSimplified++
		return expr.Right
	}

	// Check for float versions
	if rightFloat, ok := expr.Right.(*parser.FloatLiteral); ok {
		if rightFloat.Value == 0.0 {
			o.stats.AlgebraicSimplified++
			return &parser.FloatLiteral{Value: 0.0, Pos: expr.Pos}
		}
		if rightFloat.Value == 1.0 {
			o.stats.AlgebraicSimplified++
			return expr.Left
		}
	}

	if leftFloat, ok := expr.Left.(*parser.FloatLiteral); ok {
		if leftFloat.Value == 0.0 {
			o.stats.AlgebraicSimplified++
			return &parser.FloatLiteral{Value: 0.0, Pos: expr.Pos}
		}
		if leftFloat.Value == 1.0 {
			o.stats.AlgebraicSimplified++
			return expr.Right
		}
	}

	return expr
}

func (o *Optimizer) simplifyDivision(expr *parser.BinaryExpr) parser.Expr {
	// x / 1 → x
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 1 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	if rightFloat, ok := expr.Right.(*parser.FloatLiteral); ok && rightFloat.Value == 1.0 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	// 0 / x → 0 (if x != 0)
	leftVal, lok := getNumericValue(expr.Left)
	rightVal, _ := getNumericValue(expr.Right)

	if lok && leftVal == 0 && rightVal != 0 {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 0, Pos: expr.Pos}
	}

	// x / x → 1 (if both are identical and not zero)
	if exprEqual(expr.Left, expr.Right) {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 1, Pos: expr.Pos}
	}

	return expr
}

func (o *Optimizer) simplifyPower(expr *parser.BinaryExpr) parser.Expr {
	// x ^ 0 → 1
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 0 {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 1, Pos: expr.Pos}
	}

	if rightFloat, ok := expr.Right.(*parser.FloatLiteral); ok && rightFloat.Value == 0.0 {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 1, Pos: expr.Pos}
	}

	// x ^ 1 → x
	if rightInt, ok := expr.Right.(*parser.IntLiteral); ok && rightInt.Value == 1 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	if rightFloat, ok := expr.Right.(*parser.FloatLiteral); ok && rightFloat.Value == 1.0 {
		o.stats.AlgebraicSimplified++
		return expr.Left
	}

	// 0 ^ x → 0 (for x > 0)
	if leftInt, ok := expr.Left.(*parser.IntLiteral); ok && leftInt.Value == 0 {
		rightVal, rok := getNumericValue(expr.Right)
		if rok && rightVal > 0 {
			o.stats.AlgebraicSimplified++
			return &parser.IntLiteral{Value: 0, Pos: expr.Pos}
		}
	}

	// 1 ^ x → 1
	if leftInt, ok := expr.Left.(*parser.IntLiteral); ok && leftInt.Value == 1 {
		o.stats.AlgebraicSimplified++
		return &parser.IntLiteral{Value: 1, Pos: expr.Pos}
	}

	return expr
}

func (o *Optimizer) simplifyUnaryExpr(expr *parser.UnaryExpr) parser.Expr {
	// !!x → x (double negation)
	if expr.Operator == "!" {
		if innerUnary, ok := expr.Right.(*parser.UnaryExpr); ok {
			if innerUnary.Operator == "!" {
				o.stats.AlgebraicSimplified++
				return innerUnary.Right
			}
		}
	}

	// --x → x (double negative)
	if expr.Operator == "-" {
		if innerUnary, ok := expr.Right.(*parser.UnaryExpr); ok {
			if innerUnary.Operator == "-" {
				o.stats.AlgebraicSimplified++
				return innerUnary.Right
			}
		}
	}

	return expr
}
