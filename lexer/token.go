package lexer

// Type is a token type.
type Type string

const (
	ILLEGAL Type = "ILLEGAL"
	EOF     Type = "EOF"

	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	PLUS     = "+"
	MINUS    = "-"
	ASTARISK = "*"
	SLASH    = "/"
	MODULO   = "%"
	POWER    = "^"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="
	EQ  = "="
	NEQ = "!="

	AND = "&&"
	OR  = "||"
	NOT = "!"

	REGEX = "~"

	COMMA = ","

	LPAREN = "("
	RPAREN = ")"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]Type{
	"true":  TRUE,
	"false": FALSE,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
