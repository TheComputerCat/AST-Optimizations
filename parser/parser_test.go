package parser

import (
	"express_compiler/lexer"
	"testing"
)

func TestIntegerLiterals(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l)
	expr := p.ParseExpression()

	checkParserErrors(t, p)

	literal, ok := expr.(*IntLiteral)
	if !ok {
		t.Fatalf("expr not *IntLiteral. got=%T", expr)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
}

func TestFloatLiterals(t *testing.T) {
	input := "3.14"

	l := lexer.New(input)
	p := New(l)
	expr := p.ParseExpression()

	checkParserErrors(t, p)

	literal, ok := expr.(*FloatLiteral)
	if !ok {
		t.Fatalf("expr not *FloatLiteral. got=%T", expr)
	}

	if literal.Value != 3.14 {
		t.Errorf("literal.Value not %f. got=%f", 3.14, literal.Value)
	}
}

func TestStringLiterals(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	expr := p.ParseExpression()

	checkParserErrors(t, p)

	literal, ok := expr.(*StringLiteral)
	if !ok {
		t.Fatalf("expr not *StringLiteral. got=%T", expr)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestBooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		expr := p.ParseExpression()

		checkParserErrors(t, p)

		boolean, ok := expr.(*BoolLiteral)
		if !ok {
			t.Fatalf("expr not *BoolLiteral. got=%T", expr)
		}

		if boolean.Value != tt.expected {
			t.Errorf("boolean.Value not %t. got=%t", tt.expected, boolean.Value)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	expr := p.ParseExpression()

	checkParserErrors(t, p)

	ident, ok := expr.(*Identifier)
	if !ok {
		t.Fatalf("expr not *Identifier. got=%T", expr)
	}

	if ident.Name != "foobar" {
		t.Errorf("ident.Name not %s. got=%s", "foobar", ident.Name)
	}
}

func TestUnaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		expr := p.ParseExpression()

		checkParserErrors(t, p)

		unary, ok := expr.(*UnaryExpr)
		if !ok {
			t.Fatalf("expr not *UnaryExpr. got=%T", expr)
		}

		if unary.Operator != tt.operator {
			t.Fatalf("unary.Operator is not '%s'. got=%s",
				tt.operator, unary.Operator)
		}

		if !testLiteralExpression(t, unary.Right, tt.value) {
			return
		}
	}
}

func TestBinaryExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 % 2", 5, "%", 2},
		{"5 ^ 2", 5, "^", 2},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true && false", true, "&&", false},
		{"true || false", true, "||", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		expr := p.ParseExpression()

		checkParserErrors(t, p)

		if !testBinaryExpression(t, expr, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2 + 3", "(+ 2 3)"},
		{"2 + 3 * 4", "(+ 2 (* 3 4))"},
		{"2 * 3 + 4", "(+ (* 2 3) 4)"},
		{"(2 + 3) * 4", "(* (+ 2 3) 4)"},
		{"2 + 3 + 4", "(+ (+ 2 3) 4)"},
		{"2 + 3 - 4", "(- (+ 2 3) 4)"},
		{"2 * 3 * 4", "(* (* 2 3) 4)"},
		{"2 * 3 / 4", "(/ (* 2 3) 4)"},
		{"-5 + 3", "(+ (-5) 3)"},
		{"!true", "(!true)"},
		{"!(x > 5)", "(!(> x 5))"},
		{"x > 5 && y < 10", "(&& (> x 5) (< y 10))"},
		{"x > 5 || y < 10", "(|| (> x 5) (< y 10))"},
		{"2 ^ 3 ^ 2", "(^ 2 (^ 3 2))"}, // Right-associative
		{"2 + 3 ^ 2", "(+ 2 (^ 3 2))"}, // Power has higher precedence
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		expr := p.ParseExpression()

		checkParserErrors(t, p)

		actual := expr.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestGroupedExpressions(t *testing.T) {
	input := "(5 + 5) * 2"

	l := lexer.New(input)
	p := New(l)
	expr := p.ParseExpression()

	checkParserErrors(t, p)

	binary, ok := expr.(*BinaryExpr)
	if !ok {
		t.Fatalf("expr not *BinaryExpr. got=%T", expr)
	}

	if binary.Operator != "*" {
		t.Fatalf("binary.Operator is not '*'. got=%s", binary.Operator)
	}

	// Left side should be a grouped expression (another BinaryExpr)
	leftBinary, ok := binary.Left.(*BinaryExpr)
	if !ok {
		t.Fatalf("binary.Left not *BinaryExpr. got=%T", binary.Left)
	}

	if leftBinary.Operator != "+" {
		t.Fatalf("leftBinary.Operator is not '+'. got=%s", leftBinary.Operator)
	}
}

func TestCallExpressions(t *testing.T) {
	tests := []struct {
		input        string
		functionName string
		numArgs      int
	}{
		{"max(1, 2)", "max", 2},
		{"min(x, y)", "min", 2},
		{"sqrt(16)", "sqrt", 1},
		{"foo()", "foo", 0},
		{"bar(1, 2, 3)", "bar", 3},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		expr := p.ParseExpression()

		checkParserErrors(t, p)

		call, ok := expr.(*CallExpr)
		if !ok {
			t.Fatalf("expr not *CallExpr. got=%T", expr)
		}

		if call.Function != tt.functionName {
			t.Errorf("call.Function not %s. got=%s",
				tt.functionName, call.Function)
		}

		if len(call.Args) != tt.numArgs {
			t.Errorf("wrong number of arguments. want=%d, got=%d",
				tt.numArgs, len(call.Args))
		}
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"x * 2 + y * 3",
			"(+ (* x 2) (* y 3))",
		},
		{
			"(x + y) * (a + b)",
			"(* (+ x y) (+ a b))",
		},
		{
			"max(x + 1, y - 2)",
			"max((+ x 1), (- y 2))",
		},
		{
			"!flag && x > 0",
			"(&& (!flag) (> x 0))",
		},
		{
			`"hello" + " " + "world"`,
			`(+ (+ "hello" " ") "world")`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		expr := p.ParseExpression()

		checkParserErrors(t, p)

		actual := expr.String()
		if actual != tt.expected {
			t.Errorf("for input %q:\nexpected=%q\ngot=%q",
				tt.input, tt.expected, actual)
		}
	}
}

// Test error cases
func TestParserErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{
			"(2 + 3",
			"expected next token to be )",
		},
		{
			"max(1, 2",
			"expected next token to be )",
		},
		{
			"@",
			"no prefix parse function for ILLEGAL",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		p.ParseExpression()

		if len(p.Errors()) == 0 {
			t.Errorf("expected parser errors for input %q, got none", tt.input)
			continue
		}

		// Check that error message contains expected substring
		found := false
		for _, err := range p.Errors() {
			if contains(err, tt.expectedErr) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected error containing %q, got errors: %v",
				tt.expectedErr, p.Errors())
		}
	}
}

// Helper functions

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLiteralExpression(t *testing.T, expr Expr, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntLiteral(t, expr, int64(v))
	case int64:
		return testIntLiteral(t, expr, v)
	case float64:
		return testFloatLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBoolLiteral(t, expr, v)
	}
	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}

func testIntLiteral(t *testing.T, expr Expr, value int64) bool {
	integ, ok := expr.(*IntLiteral)
	if !ok {
		t.Errorf("expr not *IntLiteral. got=%T", expr)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, expr Expr, value float64) bool {
	floatLit, ok := expr.(*FloatLiteral)
	if !ok {
		t.Errorf("expr not *FloatLiteral. got=%T", expr)
		return false
	}

	if floatLit.Value != value {
		t.Errorf("floatLit.Value not %f. got=%f", value, floatLit.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expr Expr, value string) bool {
	ident, ok := expr.(*Identifier)
	if !ok {
		t.Errorf("expr not *Identifier. got=%T", expr)
		return false
	}

	if ident.Name != value {
		t.Errorf("ident.Name not %s. got=%s", value, ident.Name)
		return false
	}

	return true
}

func testBoolLiteral(t *testing.T, expr Expr, value bool) bool {
	boolean, ok := expr.(*BoolLiteral)
	if !ok {
		t.Errorf("expr not *BoolLiteral. got=%T", expr)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	return true
}

func testBinaryExpression(t *testing.T, expr Expr, left interface{},
	operator string, right interface{}) bool {

	opExpr, ok := expr.(*BinaryExpr)
	if !ok {
		t.Errorf("expr is not *BinaryExpr. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("expr.Operator is not '%s'. got=%q", operator, opExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}

	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || contains(s[1:], substr)))
}
