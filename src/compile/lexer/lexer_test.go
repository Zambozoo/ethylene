package lexer

import (
	"geth-cody/compile/lexer/token"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name     string
	input    string
	expected []token.Token
	errFunc  assert.ErrorAssertionFunc
}

func testLexHelper(t *testing.T, testCases []testCase) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := NewLexer(tt.input, nil).Lex()

			tt.errFunc(t, err)

			if !reflect.DeepEqual(tokens, tt.expected) {
				t.Error(cmp.Diff(tokens, tt.expected))
			}
		})
	}
}

func Test_Lex(t *testing.T) {
	testCases := []testCase{
		{
			name:  "empty input",
			input: "",
			expected: []token.Token{
				{Type: token.TOK_EOF, Loc: token.Location{}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "simple class",
			input: "class Class {}",
			expected: []token.Token{
				{Type: token.TOK_CLASS, Loc: token.Location{EndColumn: 5}},
				{Type: token.TOK_IDENTIFIER, Value: "Class", Loc: token.Location{StartColumn: 6, EndColumn: 11}},
				{Type: token.TOK_LEFTBRACE, Loc: token.Location{StartColumn: 12, EndColumn: 13}},
				{Type: token.TOK_RIGHTBRACE, Loc: token.Location{StartColumn: 13, EndColumn: 14}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 14, EndColumn: 14}},
			},
			errFunc: assert.NoError,
		},
	}

	testLexHelper(t, testCases)
}

func Test_LexChars(t *testing.T) {
	testCases := []testCase{
		{
			name:  "simple character",
			input: "'a'",
			expected: []token.Token{
				{Type: token.TOK_CHARACTER, Rune: 'a', Value: "a", Loc: token.Location{EndColumn: 3}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 3, EndColumn: 3}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "uncode character",
			input: "'√'",
			expected: []token.Token{
				{Type: token.TOK_CHARACTER, Rune: '√', Value: "√", Loc: token.Location{EndColumn: 3}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 3, EndColumn: 3}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "escaped uncode character",
			input: `'\u221a'`,
			expected: []token.Token{
				{Type: token.TOK_CHARACTER, Rune: '√', Value: "√", Loc: token.Location{EndColumn: 8}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 8, EndColumn: 8}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "unfinished character",
			input: "'a",
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: "'a", Loc: token.Location{EndColumn: 2}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 2, EndColumn: 2}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "unfinished unicode character",
			input: `'\u`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: `'\u`, Loc: token.Location{EndColumn: 3}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 3, EndColumn: 3}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "invalid unicode character, terminated",
			input: `'\uwxyz'`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: `'\uwxyz'`, Loc: token.Location{EndColumn: 8}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 8, EndColumn: 8}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "invalid unicode character, unterminated",
			input: `'\uwxyz`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: `'\uwxyz`, Loc: token.Location{EndColumn: 7}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 7, EndColumn: 7}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "newline in character",
			input: "'\n'",
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: "'\n'", Loc: token.Location{EndLine: 1, EndColumn: 1}},
				{Type: token.TOK_EOF, Loc: token.Location{StartLine: 1, StartColumn: 1, EndLine: 1, EndColumn: 1}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "newline in unterminated character",
			input: "'\n",
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: "'\n", Loc: token.Location{EndLine: 1}},
				{Type: token.TOK_EOF, Loc: token.Location{StartLine: 1, EndLine: 1}},
			},
			errFunc: assert.Error,
		},
	}

	testLexHelper(t, testCases)
}

func Test_LexString(t *testing.T) {
	testCases := []testCase{
		{
			name:  "simple string",
			input: `"asdf"`,
			expected: []token.Token{
				{Type: token.TOK_STRING, Value: "asdf", Loc: token.Location{EndColumn: 6}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 6, EndColumn: 6}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "uncode string",
			input: `"√"`,
			expected: []token.Token{
				{Type: token.TOK_STRING, Value: "√", Loc: token.Location{EndColumn: 3}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 3, EndColumn: 3}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "escaped uncode string",
			input: `"\u221a"`,
			expected: []token.Token{
				{Type: token.TOK_STRING, Value: "√", Loc: token.Location{EndColumn: 8}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 8, EndColumn: 8}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "invalid unicode character, terminated",
			input: `"\uwxyz"`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: `"\uwxyz"`, Loc: token.Location{EndColumn: 8}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 8, EndColumn: 8}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "invalid unicode character, unterminated",
			input: `"\uwxyz`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: `"\uwxyz`, Loc: token.Location{EndColumn: 7}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 7, EndColumn: 7}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "unfinished unicode character",
			input: `"\u`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: `"\u`, Loc: token.Location{EndColumn: 3}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 3, EndColumn: 3}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "newline in string",
			input: "\"\n\"",
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: "\"\n\"", Loc: token.Location{EndLine: 1, EndColumn: 1}},
				{Type: token.TOK_EOF, Loc: token.Location{StartLine: 1, StartColumn: 1, EndLine: 1, EndColumn: 1}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "newline in unterminated string",
			input: "\"\n",
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: "\"\n", Loc: token.Location{EndLine: 1}},
				{Type: token.TOK_EOF, Loc: token.Location{StartLine: 1, EndLine: 1}},
			},
			errFunc: assert.Error,
		},
	}

	testLexHelper(t, testCases)
}

func Test_LexInt(t *testing.T) {
	testCases := []testCase{
		{
			name:  "decimal",
			input: `10`,
			expected: []token.Token{
				{Type: token.TOK_INTEGER, Integer: 10, Value: "10", Loc: token.Location{EndColumn: 2}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 2, EndColumn: 2}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "hexadecimal",
			input: `0x10`,
			expected: []token.Token{
				{Type: token.TOK_INTEGER, Integer: 16, Value: "0x10", Loc: token.Location{EndColumn: 4}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 4, EndColumn: 4}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "invalid hexadecimal",
			input: `0xfffffffffffffffff`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: "0xfffffffffffffffff", Loc: token.Location{EndColumn: 19}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 19, EndColumn: 19}},
			},
			errFunc: assert.Error,
		},
		{
			name:  "binary",
			input: `0b10`,
			expected: []token.Token{
				{Type: token.TOK_INTEGER, Integer: 2, Value: "0b10", Loc: token.Location{EndColumn: 4}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 4, EndColumn: 4}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "invalid binary",
			input: `0b11111111111111111111111111111111111111111111111111111111111111111`,
			expected: []token.Token{
				{Type: token.TOK_UNKOWN, Value: "0b11111111111111111111111111111111111111111111111111111111111111111", Loc: token.Location{EndColumn: 67}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 67, EndColumn: 67}},
			},
			errFunc: assert.Error,
		},
	}

	testLexHelper(t, testCases)
}

func Test_LexFloat(t *testing.T) {
	testCases := []testCase{
		{
			name:  "decimal",
			input: `10e10`,
			expected: []token.Token{
				{Type: token.TOK_FLOAT, Float: 10e10, Value: "10e10", Loc: token.Location{EndColumn: 5}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 5, EndColumn: 5}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "decimal with negative exponent",
			input: `10e-10`,
			expected: []token.Token{
				{Type: token.TOK_FLOAT, Float: 10e-10, Value: "10e-10", Loc: token.Location{EndColumn: 6}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 6, EndColumn: 6}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "hexadecimal",
			input: `0x10p10`,
			expected: []token.Token{
				{Type: token.TOK_FLOAT, Float: 0x10p10, Value: "0x10p10", Loc: token.Location{EndColumn: 7}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 7, EndColumn: 7}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "hexadecimal wth negative exponent",
			input: `0x10p-10`,
			expected: []token.Token{
				{Type: token.TOK_FLOAT, Float: 0x10p-10, Value: "0x10p-10", Loc: token.Location{EndColumn: 8}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 8, EndColumn: 8}},
			},
			errFunc: assert.NoError,
		},
	}

	testLexHelper(t, testCases)
}

func Test_LexIdentifiers(t *testing.T) {
	testCases := []testCase{
		{
			name:  "simple identifier",
			input: `asdf`,
			expected: []token.Token{
				{Type: token.TOK_IDENTIFIER, Value: "asdf", Loc: token.Location{EndColumn: 4}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 4, EndColumn: 4}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "underscore-prefixed identifier",
			input: `_asdf`,
			expected: []token.Token{
				{Type: token.TOK_IDENTIFIER, Value: "_asdf", Loc: token.Location{EndColumn: 5}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 5, EndColumn: 5}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "underscore",
			input: `_`,
			expected: []token.Token{
				{Type: token.TOK_IDENTIFIER, Value: "_", Loc: token.Location{EndColumn: 1}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 1, EndColumn: 1}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "numbered identifier",
			input: `a5`,
			expected: []token.Token{
				{Type: token.TOK_IDENTIFIER, Value: "a5", Loc: token.Location{EndColumn: 2}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 2, EndColumn: 2}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "snakecase identifier",
			input: `i_am_a_snake`,
			expected: []token.Token{
				{Type: token.TOK_IDENTIFIER, Value: "i_am_a_snake", Loc: token.Location{EndColumn: 12}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 12, EndColumn: 12}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "camelcase identifier",
			input: `iAmACamel`,
			expected: []token.Token{
				{Type: token.TOK_IDENTIFIER, Value: "iAmACamel", Loc: token.Location{EndColumn: 9}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 9, EndColumn: 9}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "titlecase identifier",
			input: `IAmATitle`,
			expected: []token.Token{
				{Type: token.TOK_IDENTIFIER, Value: "IAmATitle", Loc: token.Location{EndColumn: 9}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: 9, EndColumn: 9}},
			},
			errFunc: assert.NoError,
		},
	}

	testLexHelper(t, testCases)
}

func Test_LexKeywords(t *testing.T) {
	var testCases []testCase
	for keyword := range token.KeywordMap {
		length := len(keyword)
		testCases = append(testCases, testCase{
			name:  string(keyword),
			input: string(keyword),
			expected: []token.Token{
				{Type: keyword, Loc: token.Location{EndColumn: length}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: length, EndColumn: length}},
			},
			errFunc: assert.NoError,
		})
	}

	testLexHelper(t, testCases)
}

func Test_LexSymbols(t *testing.T) {
	var testCases []testCase
	for symbol := range token.SymbolMap {
		length := len(symbol)
		testCases = append(testCases, testCase{
			name:  string(symbol),
			input: string(symbol),
			expected: []token.Token{
				{Type: symbol, Loc: token.Location{EndColumn: length}},
				{Type: token.TOK_EOF, Loc: token.Location{StartColumn: length, EndColumn: length}},
			},
			errFunc: assert.NoError,
		})
	}

	testLexHelper(t, testCases)
}
