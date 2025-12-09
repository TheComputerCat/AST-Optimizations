package main

import (
	"express_compiler/evaluator"
	"express_compiler/lexer"
	"express_compiler/optimizer"
	"express_compiler/parser"
	"fmt"
)

func main() {
	tests := []struct {
		name      string
		input     string
		constants map[string]any
		variables map[string]evaluator.Type // For type checking
	}{
		{
			name:  "Valid: Constant folding",
			input: "2 + 3 * 4",
		},
		{
			name:  "Valid: Constant propagation",
			input: "x * 2 + y",
			constants: map[string]any{
				"x": 10,
			},
			variables: map[string]evaluator.Type{
				"x": evaluator.TypeInt,
				"y": evaluator.TypeInt,
			},
		},
		{
			name:  "Valid: Algebraic simplification",
			input: "x * 0 + y",
			constants: map[string]any{
				"y": 5,
			},
			variables: map[string]evaluator.Type{
				"x": evaluator.TypeInt,
				"y": evaluator.TypeInt,
			},
		},
		{
			name:  "Valid: Strength reduction",
			input: "x * 2 + y ^ 2",
			variables: map[string]evaluator.Type{
				"x": evaluator.TypeInt,
				"y": evaluator.TypeInt,
			},
		},
		{
			name:  "Valid: Short-circuit",
			input: "true || (x > 5)",
			variables: map[string]evaluator.Type{
				"x": evaluator.TypeInt,
			},
		},
		{
			name:  "Invalid: regex",
			input: `"test" ~ "[invalid(regex"`,
		},
		{
			name:  "Valid: regex",
			input: `"test@example.com" ~ ("^[a-z]+@[a-z]" + "+\.[a-z]+$")`,
		},
		{
			name:  "Valid: regex",
			input: `"test@example.com" ~ ("^[a-z]+@[a-z]" + x)`,
			constants: map[string]any{
				"x": `"+\.[a-z]+$"`,
			},
			variables: map[string]evaluator.Type{
				"x": evaluator.TypeString,
			},
		},
		{
			name:  "Function calls",
			input: "max(2 + 1, 5 - 2)",
		},
		{
			name:  "Invalid: Type mismatch",
			input: `"hello" + 5`,
		},
		{
			name:  "Invalid: Wrong operator",
			input: `true * 3`,
		},
		{
			name:  "Valid: String concatenation",
			input: `"hello" + " " + "world"`,
		},
		{
			name:  "Invalid: Logical on numbers",
			input: `5 && 10`,
		},
		{
			name:  "Valid: Boolean logic",
			input: `x > 5 && y < 10`,
			variables: map[string]evaluator.Type{
				"x": evaluator.TypeInt,
				"y": evaluator.TypeInt,
			},
		},
		{
			name:  "Invalid: Cannot negate string",
			input: `-"hello"`,
		},
		{
			name:  "Invalid: Wrong function arity",
			input: `max(1, 2, 3)`,
		},
		{
			name:  "Valid: Function call",
			input: `max(x, y) + sqrt(16)`,
			variables: map[string]evaluator.Type{
				"x": evaluator.TypeInt,
				"y": evaluator.TypeInt,
			},
		},
	}

	for _, tt := range tests {
		fmt.Printf("\n" + "==========================================" + "\n")
		fmt.Printf("Test: %s\n", tt.name)
		fmt.Printf("Input: %s\n", tt.input)

		l := lexer.New(tt.input)

		p := parser.New(l)
		expr := p.ParseExpression()

		if len(p.Errors()) > 0 {
			fmt.Println("Parser errors:")
			for _, err := range p.Errors() {
				fmt.Printf("  %s\n", err)
			}
			continue
		}

		fmt.Printf("\nAST: %s\n", expr.String())

		tc := evaluator.NewTypeChecker(map[string]evaluator.Type{})

		if tt.variables != nil {
			tc.WithVariables(tt.variables)
		}

		exprType := tc.Check(expr)

		if tc.HasErrors() {
			fmt.Println("Type errors:")
			for _, err := range tc.Errors() {
				fmt.Printf("  %s\n", err)
			}
			continue
		}

		fmt.Printf("Type: %s\n", exprType)

		var opt *optimizer.Optimizer
		if len(tt.constants) > 0 {
			opt = optimizer.NewWithConstants(tt.constants)
		} else {
			opt = optimizer.New()
		}

		optimized := opt.Optimize(expr)

		fmt.Printf("\nOptimized: %s\n", optimized.String())

		stats := opt.Stats()

		fmt.Printf("\nOptimizations applied:\n")
		if stats.ConstantsFolded > 0 {
			fmt.Printf("  - Constants folded: %d\n", stats.ConstantsFolded)
		}
		if stats.AlgebraicSimplified > 0 {
			fmt.Printf("  - Algebraic simplifications: %d\n", stats.AlgebraicSimplified)
		}
		if stats.StrengthReduced > 0 {
			fmt.Printf("  - Strength reductions: %d\n", stats.StrengthReduced)
		}
		if stats.ShortCircuited > 0 {
			fmt.Printf("  - Strength reductions: %d\n", stats.StrengthReduced)
		}
		if stats.RegexErrors != nil {
			fmt.Printf("\nRegex validation errors:\n")
			for _, err := range stats.RegexErrors {
				fmt.Printf("  - %s\n", err)
			}
		}
	}
}
