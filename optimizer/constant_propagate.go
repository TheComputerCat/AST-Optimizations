package optimizer

import "express_compiler/parser"

// constantPropagate replaces variables with their known constant values
func (o *Optimizer) constantPropagate(expr parser.Expr) parser.Expr {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *parser.Identifier:
		// Check if we have a constant value for this identifier
		if val, ok := o.constants[e.Name]; ok {
			o.stats.VariablesPropagated++
			return o.valueToExpr(val, e.Pos)
		}
		return e

	case *parser.BinaryExpr:
		// Recursively propagate in children
		e.Left = o.constantPropagate(e.Left)
		e.Right = o.constantPropagate(e.Right)
		return e

	case *parser.UnaryExpr:
		e.Right = o.constantPropagate(e.Right)
		return e

	case *parser.CallExpr:
		for i, arg := range e.Args {
			e.Args[i] = o.constantPropagate(arg)
		}
		return e

	// Literals don't need propagation
	default:
		return e
	}
}

// valueToExpr converts a Go value to an AST expression
func (o *Optimizer) valueToExpr(val interface{}, pos parser.Position) parser.Expr {
	switch v := val.(type) {
	case int:
		return &parser.IntLiteral{Value: int64(v), Pos: pos}
	case int64:
		return &parser.IntLiteral{Value: v, Pos: pos}
	case float64:
		return &parser.FloatLiteral{Value: v, Pos: pos}
	case string:
		return &parser.StringLiteral{Value: v, Pos: pos}
	case bool:
		return &parser.BoolLiteral{Value: v, Pos: pos}
	default:
		// Unknown type - return identifier unchanged
		return &parser.Identifier{Name: "unknown", Pos: pos}
	}
}
