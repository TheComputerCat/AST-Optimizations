package optimizer

import (
	"express_compiler/parser"
	"math"
)

// constantFold evaluates constant expressions at compile time
func (o *Optimizer) constantFold(expr parser.Expr) parser.Expr {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *parser.BinaryExpr:
		// First fold children (bottom-up)
		e.Left = o.constantFold(e.Left)
		e.Right = o.constantFold(e.Right)

		// Then try to fold this node
		return o.foldBinaryExpr(e)

	case *parser.UnaryExpr:
		// Fold child first
		e.Right = o.constantFold(e.Right)

		// Then try to fold this node
		return o.foldUnaryExpr(e)

	case *parser.CallExpr:
		// Fold arguments first
		for i, arg := range e.Args {
			e.Args[i] = o.constantFold(arg)
		}

		// Then try to fold the call
		return o.foldCallExpr(e)

	// Literals and identifiers are already folded
	default:
		return e
	}
}

func (o *Optimizer) foldBinaryExpr(expr *parser.BinaryExpr) parser.Expr {
	switch expr.Operator {
	case "+":
		return o.foldAddition(expr)
	case "-":
		return o.foldSubtraction(expr)
	case "*":
		return o.foldMultiplication(expr)
	case "/":
		return o.foldDivision(expr)
	case "%":
		return o.foldModulo(expr)
	case "^":
		return o.foldPower(expr)
	case "==", "!=", "<", ">", "<=", ">=":
		return o.foldComparison(expr)
	case "&&", "||":
		return o.foldLogical(expr)
	case "~":
		// Regex matching - can't fold at compile time (would need regex engine)
		// Could potentially fold if both operands are constants, but skip for now
		return expr
	default:
		return expr
	}
}

func (o *Optimizer) foldAddition(expr *parser.BinaryExpr) parser.Expr {
	// Try int + int
	if leftInt, lok := expr.Left.(*parser.IntLiteral); lok {
		if rightInt, rok := expr.Right.(*parser.IntLiteral); rok {
			o.stats.ConstantsFolded++
			return &parser.IntLiteral{
				Value: leftInt.Value + rightInt.Value,
				Pos:   expr.Pos,
			}
		}
	}

	// Try numeric + numeric (promote to float)
	leftVal, lok := getNumericValue(expr.Left)
	rightVal, rok := getNumericValue(expr.Right)
	if lok && rok {
		o.stats.ConstantsFolded++
		return &parser.FloatLiteral{
			Value: leftVal + rightVal,
			Pos:   expr.Pos,
		}
	}

	// Try string + string
	if leftStr, lok := expr.Left.(*parser.StringLiteral); lok {
		if rightStr, rok := expr.Right.(*parser.StringLiteral); rok {
			o.stats.ConstantsFolded++
			return &parser.StringLiteral{
				Value: leftStr.Value + rightStr.Value,
				Pos:   expr.Pos,
			}
		}
	}

	// Can't fold - return unchanged
	return expr
}

func (o *Optimizer) foldSubtraction(expr *parser.BinaryExpr) parser.Expr {
	// Try int - int
	if leftInt, lok := expr.Left.(*parser.IntLiteral); lok {
		if rightInt, rok := expr.Right.(*parser.IntLiteral); rok {
			o.stats.ConstantsFolded++
			return &parser.IntLiteral{
				Value: leftInt.Value - rightInt.Value,
				Pos:   expr.Pos,
			}
		}
	}

	// Try numeric - numeric
	leftVal, lok := getNumericValue(expr.Left)
	rightVal, rok := getNumericValue(expr.Right)
	if lok && rok {
		o.stats.ConstantsFolded++
		return &parser.FloatLiteral{
			Value: leftVal - rightVal,
			Pos:   expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) foldMultiplication(expr *parser.BinaryExpr) parser.Expr {
	// Try int * int
	if leftInt, lok := expr.Left.(*parser.IntLiteral); lok {
		if rightInt, rok := expr.Right.(*parser.IntLiteral); rok {
			o.stats.ConstantsFolded++
			return &parser.IntLiteral{
				Value: leftInt.Value * rightInt.Value,
				Pos:   expr.Pos,
			}
		}
	}

	// Try numeric * numeric
	leftVal, lok := getNumericValue(expr.Left)
	rightVal, rok := getNumericValue(expr.Right)
	if lok && rok {
		o.stats.ConstantsFolded++
		return &parser.FloatLiteral{
			Value: leftVal * rightVal,
			Pos:   expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) foldDivision(expr *parser.BinaryExpr) parser.Expr {
	leftVal, lok := getNumericValue(expr.Left)
	rightVal, rok := getNumericValue(expr.Right)

	if lok && rok {
		// Check for division by zero
		if rightVal == 0 {
			// Don't optimize - let runtime handle the error
			return expr
		}

		o.stats.ConstantsFolded++
		return &parser.FloatLiteral{
			Value: leftVal / rightVal,
			Pos:   expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) foldModulo(expr *parser.BinaryExpr) parser.Expr {
	// Modulo only works with integers
	leftInt, lok := getIntValue(expr.Left)
	rightInt, rok := getIntValue(expr.Right)

	if lok && rok {
		// Check for division by zero
		if rightInt == 0 {
			return expr
		}

		o.stats.ConstantsFolded++
		return &parser.IntLiteral{
			Value: leftInt % rightInt,
			Pos:   expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) foldPower(expr *parser.BinaryExpr) parser.Expr {
	leftVal, lok := getNumericValue(expr.Left)
	rightVal, rok := getNumericValue(expr.Right)

	if lok && rok {
		o.stats.ConstantsFolded++
		return &parser.FloatLiteral{
			Value: math.Pow(leftVal, rightVal),
			Pos:   expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) foldComparison(expr *parser.BinaryExpr) parser.Expr {
	// Try numeric comparison
	leftVal, lok := getNumericValue(expr.Left)
	rightVal, rok := getNumericValue(expr.Right)

	if lok && rok {
		var result bool
		switch expr.Operator {
		case "==":
			result = leftVal == rightVal
		case "!=":
			result = leftVal != rightVal
		case "<":
			result = leftVal < rightVal
		case ">":
			result = leftVal > rightVal
		case "<=":
			result = leftVal <= rightVal
		case ">=":
			result = leftVal >= rightVal
		}

		o.stats.ConstantsFolded++
		return &parser.BoolLiteral{
			Value: result,
			Pos:   expr.Pos,
		}
	}

	// Try boolean comparison
	leftBool, lok := expr.Left.(*parser.BoolLiteral)
	rightBool, rok := expr.Right.(*parser.BoolLiteral)

	if lok && rok {
		var result bool
		switch expr.Operator {
		case "==":
			result = leftBool.Value == rightBool.Value
		case "!=":
			result = leftBool.Value != rightBool.Value
		default:
			return expr // <, >, etc. don't work with booleans
		}

		o.stats.ConstantsFolded++
		return &parser.BoolLiteral{
			Value: result,
			Pos:   expr.Pos,
		}
	}

	// Try string comparison
	leftStr, lok := expr.Left.(*parser.StringLiteral)
	rightStr, rok := expr.Right.(*parser.StringLiteral)

	if lok && rok {
		var result bool
		switch expr.Operator {
		case "==":
			result = leftStr.Value == rightStr.Value
		case "!=":
			result = leftStr.Value != rightStr.Value
		default:
			return expr // <, > don't make sense for strings in our language
		}

		o.stats.ConstantsFolded++
		return &parser.BoolLiteral{
			Value: result,
			Pos:   expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) foldLogical(expr *parser.BinaryExpr) parser.Expr {
	leftBool, lok := expr.Left.(*parser.BoolLiteral)
	rightBool, rok := expr.Right.(*parser.BoolLiteral)

	if lok && rok {
		var result bool
		switch expr.Operator {
		case "&&":
			result = leftBool.Value && rightBool.Value
		case "||":
			result = leftBool.Value || rightBool.Value
		}

		o.stats.ConstantsFolded++
		return &parser.BoolLiteral{
			Value: result,
			Pos:   expr.Pos,
		}
	}

	return expr
}

func (o *Optimizer) foldUnaryExpr(expr *parser.UnaryExpr) parser.Expr {
	switch expr.Operator {
	case "!":
		if boolLit, ok := expr.Right.(*parser.BoolLiteral); ok {
			o.stats.ConstantsFolded++
			return &parser.BoolLiteral{
				Value: !boolLit.Value,
				Pos:   expr.Pos,
			}
		}

	case "-":
		if intLit, ok := expr.Right.(*parser.IntLiteral); ok {
			o.stats.ConstantsFolded++
			return &parser.IntLiteral{
				Value: -intLit.Value,
				Pos:   expr.Pos,
			}
		}
		if floatLit, ok := expr.Right.(*parser.FloatLiteral); ok {
			o.stats.ConstantsFolded++
			return &parser.FloatLiteral{
				Value: -floatLit.Value,
				Pos:   expr.Pos,
			}
		}
	}

	return expr
}

func (o *Optimizer) foldCallExpr(expr *parser.CallExpr) parser.Expr {
	switch expr.Function {
	case "max":
		if len(expr.Args) != 2 {
			return expr // Wrong arity
		}

		arg1, ok1 := getNumericValue(expr.Args[0])
		arg2, ok2 := getNumericValue(expr.Args[1])

		if ok1 && ok2 {
			o.stats.ConstantsFolded++
			if arg1 > arg2 {
				return expr.Args[0]
			}
			return expr.Args[1]
		}

	case "min":
		if len(expr.Args) != 2 {
			return expr
		}

		arg1, ok1 := getNumericValue(expr.Args[0])
		arg2, ok2 := getNumericValue(expr.Args[1])

		if ok1 && ok2 {
			o.stats.ConstantsFolded++
			if arg1 < arg2 {
				return expr.Args[0]
			}
			return expr.Args[1]
		}

	case "sqrt":
		if len(expr.Args) != 1 {
			return expr
		}

		arg, ok := getNumericValue(expr.Args[0])
		if ok {
			o.stats.ConstantsFolded++
			return &parser.FloatLiteral{
				Value: math.Sqrt(arg),
				Pos:   expr.Pos,
			}
		}

	case "abs":
		if len(expr.Args) != 1 {
			return expr
		}

		arg, ok := getNumericValue(expr.Args[0])
		if ok {
			o.stats.ConstantsFolded++
			return &parser.FloatLiteral{
				Value: math.Abs(arg),
				Pos:   expr.Pos,
			}
		}
	}

	return expr
}
