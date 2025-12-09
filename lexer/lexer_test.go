package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `5 + 10 - 3
2.5 * 3.14
x >= y
"hello world"
true && false
max(a, b)
x ^ 2
!(x != 5)
`

	tests := []struct {
		expectedType    Type
		expectedLiteral string
	}{
		// Line 1: 5 + 10 - 3
		{INT, "5"},
		{PLUS, "+"},
		{INT, "10"},
		{MINUS, "-"},
		{INT, "3"},

		// Line 2: 2.5 * 3.14
		{FLOAT, "2.5"},
		{ASTARISK, "*"},
		{FLOAT, "3.14"},

		// Line 3: x >= y
		{IDENT, "x"},
		{GTE, ">="},
		{IDENT, "y"},

		// Line 4: "hello world"
		{STRING, "hello world"},

		// Line 5: true && false
		{TRUE, "true"},
		{AND, "&&"},
		{FALSE, "false"},

		// Line 6: max(a, b)
		{IDENT, "max"},
		{LPAREN, "("},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{RPAREN, ")"},

		// Line 7: x ^ 2
		{IDENT, "x"},
		{POWER, "^"},
		{INT, "2"},

		// Line 8: !(x != 5)
		{NOT, "!"},
		{LPAREN, "("},
		{IDENT, "x"},
		{NEQ, "!="},
		{INT, "5"},
		{RPAREN, ")"},

		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal: %q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLineColumnTracking(t *testing.T) {
	input := `x + y
z * 2`

	tests := []struct {
		expectedType   Type
		expectedLine   int
		expectedColumn int
	}{
		{IDENT, 1, 1},    // x
		{PLUS, 1, 3},     // +
		{IDENT, 1, 5},    // y
		{IDENT, 2, 1},    // z
		{ASTARISK, 2, 3}, // *
		{INT, 2, 5},      // 2
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Line != tt.expectedLine {
			t.Errorf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Errorf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

func TestFloatNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3.14", "3.14"},
		{"0.5", "0.5"},
		{"100.001", "100.001"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != FLOAT {
			t.Fatalf("expected FLOAT, got %q", tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Fatalf("expected literal %q, got %q", tt.expected, tok.Literal)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	input := `"hello"
"hello world"
"with spaces   "
""`

	tests := []string{
		"hello",
		"hello world",
		"with spaces   ",
		"",
	}

	l := New(input)

	for i, expected := range tests {
		tok := l.NextToken()

		if tok.Type != STRING {
			t.Fatalf("tests[%d] - tokentype wrong. expected=STRING, got=%q",
				i, tok.Type)
		}

		if tok.Literal != expected {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, expected, tok.Literal)
		}
	}
}

func TestTwoCharOperators(t *testing.T) {
	input := `== != <= >= && ||`

	tests := []struct {
		expectedType    Type
		expectedLiteral string
	}{
		{EQ, "=="},
		{NEQ, "!="},
		{LTE, "<="},
		{GTE, ">="},
		{AND, "&&"},
		{OR, "||"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIllegalTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected Type
	}{
		{"=", ILLEGAL}, // Single = is illegal
		{"&", ILLEGAL}, // Single & is illegal
		{"|", ILLEGAL}, // Single | is illegal
		{"@", ILLEGAL}, // @ is not in our language
		{"#", ILLEGAL}, // # is not in our language
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != tt.expected {
			t.Errorf("input %q - expected type %q, got %q",
				tt.input, tt.expected, tok.Type)
		}
	}
}

func TestAllOperators(t *testing.T) {
	input := `+ - * / % ^ ~ < > <= >= == != && || !`

	tests := []Type{
		PLUS, MINUS, ASTARISK, SLASH, MODULO, POWER, REGEX,
		LT, GT, LTE, GTE, EQ, NEQ, AND, OR, NOT,
	}

	l := New(input)

	for i, expectedType := range tests {
		tok := l.NextToken()

		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, expectedType, tok.Type)
		}
	}
}
