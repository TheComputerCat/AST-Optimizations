package optimizer

import (
	"express_compiler/parser"
)

const (
	MaxPasses = 10
)

type Optimizer struct {
	constants map[string]interface{}
	stats     *OptimizationStats
}

func New() *Optimizer {
	return &Optimizer{
		constants: make(map[string]interface{}),
		stats:     NewStats(),
	}
}

func NewWithConstants(constants map[string]interface{}) *Optimizer {
	return &Optimizer{
		constants: constants,
		stats:     NewStats(),
	}
}

// Optimize applies all optimization passes until a fixed point is reached
func (o *Optimizer) Optimize(expr parser.Expr) parser.Expr {
	if expr == nil {
		return nil
	}

	// Run optimization passes until nothing changes
	previous := expr.String()

	for pass := 0; pass < MaxPasses; pass++ {
		o.stats.PassesRun++

		expr = o.constantFold(expr)
		if len(o.constants) > 0 {
			expr = o.constantPropagate(expr)
		}
		expr = o.constantFold(expr)
		expr = o.algebraicSimplify(expr)
		expr = o.strengthReduce(expr)
		expr = o.shortCircuit(expr)
		expr = o.constantFold(expr)
		current := expr.String()
		if current == previous {
			break
		}
		previous = current
	}

	// First, validate regex patterns (doesn't modify AST)
	o.stats.RegexErrors = ValidateRegexPatterns(expr)

	return expr
}

func (o *Optimizer) Stats() *OptimizationStats {
	return o.stats
}

func getNumericValue(expr parser.Expr) (float64, bool) {
	switch e := expr.(type) {
	case *parser.IntLiteral:
		return float64(e.Value), true
	case *parser.FloatLiteral:
		return e.Value, true
	default:
		return 0, false
	}
}

func getIntValue(expr parser.Expr) (int64, bool) {
	if intLit, ok := expr.(*parser.IntLiteral); ok {
		return intLit.Value, true
	}
	return 0, false
}

func exprEqual(a, b parser.Expr) bool {
	return a.String() == b.String()
}
