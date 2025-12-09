package optimizer

import (
	"express_compiler/parser"
	"fmt"
	"regexp"
)

// ValidateRegexPatterns checks all regex patterns in the AST for validity
// Returns a list of error messages if any patterns are invalid
func ValidateRegexPatterns(expr parser.Expr) []string {
	var errors []string

	var validate func(parser.Expr)
	validate = func(e parser.Expr) {
		if e == nil {
			return
		}

		switch node := e.(type) {
		case *parser.BinaryExpr:
			if node.Operator == "~" {

				if strLit, ok := node.Right.(*parser.StringLiteral); ok {
					// Try to compile the regex
					_, err := regexp.Compile(strLit.Value)
					if err != nil {
						errors = append(errors, fmt.Sprintf(
							"line %d:%d - invalid regex pattern '%s': %v",
							node.Pos.Line, node.Pos.Column, strLit.Value, err,
						))
					}
				} else {
					// Regex pattern is not a constant - warn but don't error
					// (could be a variable holding a pattern)
					errors = append(errors, fmt.Sprintf(
						"line %d:%d - regex pattern is not a string literal, cannot validate at compile time",
						node.Pos.Line, node.Pos.Column,
					))
				}
			}

			validate(node.Left)
			validate(node.Right)

		case *parser.UnaryExpr:
			validate(node.Right)

		case *parser.CallExpr:
			for _, arg := range node.Args {
				validate(arg)
			}

		case *parser.IntLiteral, *parser.FloatLiteral,
			*parser.StringLiteral, *parser.BoolLiteral,
			*parser.Identifier:
		}
	}

	validate(expr)
	return errors
}
