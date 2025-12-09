package optimizer

import (
	"express_compiler/parser"
	"fmt"
)

type OptimizationStats struct {
	PassesRun           int
	ConstantsFolded     int
	AlgebraicSimplified int
	StrengthReduced     int
	ShortCircuited      int
	VariablesPropagated int
	CSEOpportunities    []CSEOpportunity
	RegexErrors         []string
}

type CSEOpportunity struct {
	Expression string
	Count      int
	Locations  []parser.Position
}

func NewStats() *OptimizationStats {
	return &OptimizationStats{
		CSEOpportunities: []CSEOpportunity{},
		RegexErrors:      []string{},
	}
}

func (s *OptimizationStats) Merge(other *OptimizationStats) {
	s.ConstantsFolded += other.ConstantsFolded
	s.AlgebraicSimplified += other.AlgebraicSimplified
	s.StrengthReduced += other.StrengthReduced
	s.ShortCircuited += other.ShortCircuited
	s.VariablesPropagated += other.VariablesPropagated
}

func (s *OptimizationStats) String() string {
	result := "Optimization Statistics:\n"
	result += fmt.Sprintf("  Passes run: %d\n", s.PassesRun)
	result += fmt.Sprintf("  Constants folded: %d\n", s.ConstantsFolded)
	result += fmt.Sprintf("  Variables propagated: %d\n", s.VariablesPropagated)
	result += fmt.Sprintf("  Algebraic simplifications: %d\n", s.AlgebraicSimplified)
	result += fmt.Sprintf("  Strength reductions: %d\n", s.StrengthReduced)
	result += fmt.Sprintf("  Short-circuits: %d\n", s.ShortCircuited)
	result += fmt.Sprintf("  CSE opportunities: %d\n", len(s.CSEOpportunities))

	if len(s.RegexErrors) > 0 {
		result += fmt.Sprintln("\nRegex validation errors:")
		for _, err := range s.RegexErrors {
			result += fmt.Sprintf("  - %s\n", err)
		}
	}

	return result
}
