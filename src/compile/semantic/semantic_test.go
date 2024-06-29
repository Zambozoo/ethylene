package semantic

import (
	"geth-cody/ast"
	"geth-cody/ast/decl"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/compile/syntax"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"
	"geth-cody/io/path"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name     string
	input    string
	expected *bytecode.Bytecodes
	errFunc  assert.ErrorAssertionFunc
}

type mockChan[T any] struct{}

func (c *mockChan[T]) Send(T) {}

func wrappingDecl() ast.Declaration {
	return &decl.Class{
		BaseDecl: decl.BaseDecl{
			Name_: token.Token{
				Value: "Test",
			},
		},
	}
}

func testParseHelper(t *testing.T, testCases []testCase, syntaxFunc func(*syntax.Parser) (ast.Node, io.Error), semanticFunc func(*Parser, ast.Node) (*bytecode.Bytecodes, io.Error)) {
	defer io.InitLogger()()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			filePath := path.File("test.eth")
			tokens, err := lexer.NewLexer(tt.input, &filePath).Lex()
			if err != nil {
				t.Fatal(err)
			}
			var (
				project     path.Project
				channel     mockChan[path.Path]
				mainDirPath path.File
			)
			symbolMap := syntax.SymbolMap{Types: typeid.NewTypes()}

			syntaxParser := syntax.NewParser(tokens, &project, &filePath, &mainDirPath, &path.DefaultProvider{}, &channel, symbolMap)
			syntaxParser.WrapScope(wrappingDecl())
			node, err := syntaxFunc(syntaxParser)
			if err != nil {
				t.Fatal(err)
			}

			var file ast.File
			semanticParser := NewParser(file, symbolMap)
			bytecodes, err := semanticFunc(semanticParser, node)
			tt.errFunc(t, err)

			if !reflect.DeepEqual(bytecodes, tt.expected) {
				diff := cmp.Diff(bytecodes, tt.expected)
				if diff == "" {
					diff = "an unexported field differed"
				}

				t.Error(diff)
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []testCase{
		{
			name:     "valid file",
			input:    `class Class {}`,
			expected: &bytecode.Bytecodes{},
			errFunc:  assert.NoError,
		},
		{
			name: "valid method",
			input: `class Class {
				fun int(int,int) add = (a,b) {
					return a - b;
				}
			}`,
			expected: &bytecode.Bytecodes{},
			errFunc:  assert.NoError,
		},
		{
			name: "method return type mismatch",
			input: `class Class {
				fun void(int,int) add = (a,b) {
					return a - b;
				}
			}`,
			errFunc: assert.Error,
		},
		{
			name: "valid implementing",
			input: `interface Interface {
				virtual fun int(int,int) add;

				abstract Abstract <: Interface {
					var int i;
					fun int(int,int) add = (a,b) {
						return (2*a+2*b)/2;
					}
				}

				class Class {
					var int i;
					fun int(int,int) add = (a,b) {
						return a + b;
					}
				}
			}`,
			expected: &bytecode.Bytecodes{},
			errFunc:  assert.NoError,
		},
		{
			name: "invalid implementing, mismatched type",
			input: `interface Interface {
				virtual fun int(int,int) add;

				abstract Abstract <: Interface {
				}

				class Class <: Abstract {
					fun void(int,int) add = (a,b) {
						return a + b;
					}
				}
			}`,
			errFunc: assert.Error,
		},
		{
			name: "invalid implementing, mismatched type",
			input: `interface Interface {
				virtual fun int(int,int) add;

				abstract Abstract <: Interface {
				}

				class Class {
					fun void(int,int) add = (a,b) {
						return a + b;
					}
				}
			}`,
			errFunc: assert.Error,
		},
	}

	testParseHelper(t, testCases,
		func(p *syntax.Parser) (ast.Node, io.Error) {
			return p.Parse()
		},
		func(p *Parser, n ast.Node) (*bytecode.Bytecodes, io.Error) {
			d := n.(ast.File)
			if err := d.LinkParents(p, data.NewAsyncSet[ast.Declaration]()); err != nil {
				return nil, err
			}

			if err := d.LinkFields(p, data.NewAsyncSet[ast.Declaration]()); err != nil {
				return nil, err
			}

			if err := d.Semantic(p); err != nil {
				return nil, err
			}

			return p.bytecodes, nil
		},
	)
}

func TestParseGeneric(t *testing.T) {
	testCases := []testCase{
		{
			name:     "valid generic file",
			input:    `interface Collection[E] {}`,
			expected: &bytecode.Bytecodes{},
			errFunc:  assert.NoError,
		},
		{
			name: "valid generic file with parent",
			input: `interface Example {
				interface Collection[E] {
					virtual fun bool(E) contains;
				}
				interface List[F] <: Collection[F] {
				}
				class ArrList[G] <: List[G] {
					fun bool(G) contains = (e) {
					}
				}
				static fun void(str*) main = (args) {
					var List[int]* l;
					var ArrList[int]* a;
					l = a;
				}
			}`,
			expected: &bytecode.Bytecodes{},
			errFunc:  assert.NoError,
		},
		{
			name: "valid generic file with parent and same token generic args",
			input: `interface Example {
				interface Map[K, V] {
					virtual fun bool(K) contains;
				}
				class WeirdSet[K] <: Map[K,K] {
					fun bool(K) contains = (k) {
					}
				}
			}`,
			expected: &bytecode.Bytecodes{},
			errFunc:  assert.NoError,
		},
	}

	testParseHelper(t, testCases,
		func(p *syntax.Parser) (ast.Node, io.Error) {
			return p.Parse()
		},
		func(p *Parser, n ast.Node) (*bytecode.Bytecodes, io.Error) {
			d := n.(ast.File)
			if err := d.LinkParents(p, data.NewAsyncSet[ast.Declaration]()); err != nil {
				return nil, err
			}

			if err := d.LinkFields(p, data.NewAsyncSet[ast.Declaration]()); err != nil {
				return nil, err
			}

			if err := d.Semantic(p); err != nil {
				return nil, err
			}

			return p.bytecodes, nil
		},
	)
}
