package token

type Type string

func (t Type) String() string {
	return string(t)
}

var SymbolMap = make(map[Type]struct{})

func newSymbol(s string) Type {
	SymbolMap[Type(s)] = struct{}{}
	return Type(s)
}

var KeywordMap = make(map[Type]struct{})

func newKeyword(s string) Type {
	KeywordMap[Type(s)] = struct{}{}
	return Type(s)
}

var (
	TOK_INTEGER    Type = "INT"
	TOK_FLOAT      Type = "FLOAT"
	TOK_STRING     Type = "STRING"
	TOK_CHARACTER  Type = "CHAR"
	TOK_IDENTIFIER Type = "IDENTIFIER"

	TOK_INC = newSymbol("++")
	TOK_DEC = newSymbol("--")

	TOK_PLUS   = newSymbol("+")
	TOK_MINUS  = newSymbol("-")
	TOK_STAR   = newSymbol("*")
	TOK_DIVIDE = newSymbol("/")
	TOK_MODULO = newSymbol("%")

	TOK_AND                = newSymbol("&&")
	TOK_OR                 = newSymbol("||")
	TOK_BITAND             = newSymbol("&")
	TOK_BITOR              = newSymbol("|")
	TOK_BITXOR             = newSymbol("^")
	TOK_BANG               = newSymbol("!")
	TOK_SHIFTRIGHT         = newSymbol(">>")
	TOK_SHIFTLEFT          = newSymbol("<<")
	TOK_SHIFTUNSIGNEDRIGHT = newSymbol(">>>")

	TOK_LESSTHAN         = newSymbol("<")
	TOK_LESSTHANEQUAL    = newSymbol("<=")
	TOK_GREATERTHAN      = newSymbol(">")
	TOK_GREATERTHANEQUAL = newSymbol(">=")
	TOK_SPACESHIP        = newSymbol("<=>")
	TOK_EQUAL            = newSymbol("==")
	TOK_BANGEQUAL        = newSymbol("!=")
	TOK_ASSIGN           = newSymbol("=")
	TOK_NULLARY          = newSymbol("??")

	TOK_QUESTIONMARK = newSymbol("?")
	TOK_COLON        = newSymbol(":")

	TOK_PERIOD    = newSymbol(".")
	TOK_SUBTYPE   = newSymbol("<:")
	TOK_SUPERTYPE = newSymbol(":>")
	TOK_HASHTAG   = newSymbol("#")
	TOK_AT        = newSymbol("@")
	TOK_DOLLAR    = newSymbol("$")
	TOK_COMMA     = newSymbol(",")
	TOK_SEMICOLON = newSymbol(";")
	TOK_TILDE     = newSymbol("~")

	TOK_LEFTPAREN    = newSymbol("(")
	TOK_RIGHTPAREN   = newSymbol(")")
	TOK_LEFTBRACKET  = newSymbol("[")
	TOK_RIGHTBRACKET = newSymbol("]")
	TOK_LEFTBRACE    = newSymbol("{")
	TOK_RIGHTBRACE   = newSymbol("}")

	TOK_IF       = newKeyword("if")
	TOK_ELSE     = newKeyword("else")
	TOK_FOR      = newKeyword("for")
	TOK_ASYNC    = newKeyword("async")
	TOK_WAIT     = newKeyword("wait")
	TOK_RETURN   = newKeyword("return")
	TOK_LABEL    = newKeyword("label")
	TOK_BREAK    = newKeyword("break")
	TOK_CONTINUE = newKeyword("continue")
	TOK_DELETE   = newKeyword("delete")
	TOK_FUN      = newKeyword("fun")
	TOK_VAR      = newKeyword("var")
	TOK_TYPE     = newKeyword("type")
	TOK_NEW      = newKeyword("new")
	TOK_LAMBDA   = newKeyword("lambda")

	TOK_PANIC = newKeyword("panic")
	TOK_PRINT = newKeyword("print")

	TOK_TRUE  = newKeyword("true")
	TOK_FALSE = newKeyword("false")
	TOK_THIS  = newKeyword("this")
	TOK_SUPER = newKeyword("super")
	TOK_NULL  = newKeyword("null")

	TOK_CLASS     = newKeyword("class")
	TOK_ABSTRACT  = newKeyword("abstract")
	TOK_INTERFACE = newKeyword("interface")
	TOK_STRUCT    = newKeyword("struct")
	TOK_ENUM      = newKeyword("enum")

	TOK_PUBLIC    = newKeyword("public")
	TOK_PRIVATE   = newKeyword("private")
	TOK_PROTECTED = newKeyword("protected")
	TOK_STATIC    = newKeyword("static")
	TOK_NATIVE    = newKeyword("native")
	TOK_VIRTUAL   = newKeyword("virtual")

	TOK_IMPORT = newKeyword("import")

	TOK_TYPEINT  = newKeyword("int")
	TOK_TYPEFLT  = newKeyword("flt")
	TOK_TYPECHAR = newKeyword("char")
	TOK_TYPESTR  = newKeyword("str")
	TOK_TYPEBOOL = newKeyword("bool")
	TOK_TYPEVOID = newKeyword("void")
	TOK_TYPEWORD = newKeyword("word")

	TOK_EOF    Type = "EOF"
	TOK_UNKOWN Type = "UNKNOWN"
)
