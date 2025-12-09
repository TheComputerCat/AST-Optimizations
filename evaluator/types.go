package evaluator

import (
	"express_compiler/parser"
	"fmt"
)

// Type represents the type of a value
type Type int

const (
	TypeError Type = iota
	TypeInt
	TypeFloat
	TypeString
	TypeBool
)

func (t Type) String() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	case TypeBool:
		return "bool"
	default:
		return "error"
	}
}

// TypeChecker performs type checking on expressions
type TypeChecker struct {
	variables map[string]Type // Known variable types
	errors    []string
}

func NewTypeChecker(variables map[string]Type) *TypeChecker {
	return &TypeChecker{
		variables: variables,
		errors:    []string{},
	}
}

func (tc *TypeChecker) Check(expr parser.Expr) Type {
	if expr == nil {
		return TypeError
	}

	switch e := expr.(type) {
	case *parser.IntLiteral:
		return TypeInt

	case *parser.FloatLiteral:
		return TypeFloat

	case *parser.StringLiteral:
		return TypeString

	case *parser.BoolLiteral:
		return TypeBool

	case *parser.Identifier:
		if typ, ok := tc.variables[e.Name]; ok {
			return typ
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("undefined variable: %s", e.Name))
		return TypeError

	case *parser.BinaryExpr:
		return tc.checkBinaryExpr(e)

	case *parser.UnaryExpr:
		return tc.checkUnaryExpr(e)

	case *parser.CallExpr:
		return tc.checkCallExpr(e)

	default:
		return TypeError
	}
}

func (tc *TypeChecker) checkBinaryExpr(e *parser.BinaryExpr) Type {
	leftType := tc.Check(e.Left)
	rightType := tc.Check(e.Right)

	switch e.Operator {
	case "+":
		// int + int → int
		if leftType == TypeInt && rightType == TypeInt {
			return TypeInt
		}
		// float + float → float (or int + float, etc.)
		if isNumeric(leftType) && isNumeric(rightType) {
			return TypeFloat
		}
		// string + string → string
		if leftType == TypeString && rightType == TypeString {
			return TypeString
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("cannot add %s + %s", leftType, rightType))
		return TypeError

	case "-", "*", "/", "%", "^":
		if isNumeric(leftType) && isNumeric(rightType) {
			if leftType == TypeInt && rightType == TypeInt {
				return TypeInt
			}
			return TypeFloat
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("cannot apply %s to %s and %s",
				e.Operator, leftType, rightType))
		return TypeError

	case "==", "!=", "<", ">", "<=", ">=":
		// Numeric comparison
		if isNumeric(leftType) && isNumeric(rightType) {
			return TypeBool
		}
		// String/bool equality
		if (leftType == TypeString && rightType == TypeString) ||
			(leftType == TypeBool && rightType == TypeBool) {
			return TypeBool
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("cannot compare %s and %s", leftType, rightType))
		return TypeError

	case "&&", "||":
		if leftType == TypeBool && rightType == TypeBool {
			return TypeBool
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("logical operators require bool, got %s and %s",
				leftType, rightType))
		return TypeError

	case "~":
		// Regex: string ~ string → bool
		if leftType == TypeString && rightType == TypeString {
			return TypeBool
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("regex requires strings, got %s ~ %s",
				leftType, rightType))
		return TypeError

	default:
		return TypeError
	}
}

func (tc *TypeChecker) checkUnaryExpr(e *parser.UnaryExpr) Type {
	rightType := tc.Check(e.Right)

	switch e.Operator {
	case "!":
		if rightType == TypeBool {
			return TypeBool
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("cannot negate non-boolean %s", rightType))
		return TypeError

	case "-":
		if isNumeric(rightType) {
			return rightType
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("cannot negate non-numeric %s", rightType))
		return TypeError

	default:
		return TypeError
	}
}

func (tc *TypeChecker) checkCallExpr(e *parser.CallExpr) Type {
	switch e.Function {
	case "max", "min":
		if len(e.Args) != 2 {
			tc.errors = append(tc.errors,
				fmt.Sprintf("%s requires 2 arguments, got %d",
					e.Function, len(e.Args)))
			return TypeError
		}
		arg1Type := tc.Check(e.Args[0])
		arg2Type := tc.Check(e.Args[1])
		if isNumeric(arg1Type) && isNumeric(arg2Type) {
			if arg1Type == TypeInt && arg2Type == TypeInt {
				return TypeInt
			}
			return TypeFloat
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("%s requires numeric arguments", e.Function))
		return TypeError

	case "sqrt", "abs":
		if len(e.Args) != 1 {
			tc.errors = append(tc.errors,
				fmt.Sprintf("%s requires 1 argument, got %d",
					e.Function, len(e.Args)))
			return TypeError
		}
		argType := tc.Check(e.Args[0])
		if isNumeric(argType) {
			return TypeFloat
		}
		tc.errors = append(tc.errors,
			fmt.Sprintf("%s requires numeric argument", e.Function))
		return TypeError

	default:
		tc.errors = append(tc.errors,
			fmt.Sprintf("unknown function: %s", e.Function))
		return TypeError
	}
}

func (tc *TypeChecker) Errors() []string {
	return tc.errors
}

func isNumeric(t Type) bool {
	return t == TypeInt || t == TypeFloat
}

func (tc *TypeChecker) WithVariables(vars map[string]Type) {
	for k, v := range vars {
		tc.variables[k] = v
	}
}

func (tc *TypeChecker) HasErrors() bool {
	return len(tc.errors) > 0
}
