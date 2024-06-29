package syntax

import (
	"geth-cody/ast"
	"geth-cody/ast/decl"
	"geth-cody/ast/decl/generics"
	"geth-cody/ast/expr"
	"geth-cody/ast/field"
	"geth-cody/ast/stmt"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"
	"geth-cody/io/path"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func testFile() path.Path {
	f := path.File("test.eth")
	return &f
}

type testCase struct {
	name     string
	input    string
	expected ast.Node
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

func testParseHelper(t *testing.T, testCases []testCase, f func(*Parser) (ast.Node, io.Error)) {
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			filePath := testFile()
			tokens, err := lexer.NewLexer(tt.input, filePath).Lex()
			if err != nil {
				t.Fatal(err)
			}

			var (
				project     path.Project
				channel     mockChan[path.Path]
				mainDirPath path.File
			)
			symbolMap := SymbolMap{Types: typeid.NewTypes()}

			parser := NewParser(tokens, &project, filePath, &mainDirPath, &path.DefaultProvider{}, &channel, symbolMap)
			parser.WrapScope(wrappingDecl())

			node, err := f(parser)
			tt.errFunc(t, err)

			if !reflect.DeepEqual(node, tt.expected) {
				diff := cmp.Diff(node, tt.expected, cmpopts.IgnoreUnexported(path.Project{}, typeid.Types{}))
				if diff == "" {
					diff = "an unexported field differed"
				}

				t.Error(diff)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name:    "empty string",
			input:   ``,
			errFunc: assert.Error,
		},
	}
	testParseHelper(t, testCases,
		func(p *Parser) (ast.Node, io.Error) {
			return p.Parse()
		},
	)
}

func TestParseDecl(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name:  "valid empty class",
			input: `class Class {}`,
			expected: &decl.Class{
				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_CLASS,
						Loc:  token.Location{EndColumn: 5, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 13, EndColumn: 14, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Class",
						Loc:   token.Location{StartColumn: 6, EndColumn: 11, Path_: testFile()},
					},
					IsClass:        true,
					Members_:       map[string]ast.Member{},
					Methods_:       map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty tailed class",
			input: `class Class~ {}`,
			expected: &decl.Class{
				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_CLASS,
						Loc:  token.Location{EndColumn: 5, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 14, EndColumn: 15, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Class",
						Loc:   token.Location{StartColumn: 6, EndColumn: 11, Path_: testFile()},
					},
					IsClass:        true,
					Members_:       map[string]ast.Member{},
					Methods_:       map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
				IsTailed_: true,
			},
			errFunc: assert.NoError,
		},
		func() testCase {
			c := &decl.Class{

				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_CLASS,
						Loc:  token.Location{EndColumn: 5, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 16, EndColumn: 17, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Class",
						Loc:   token.Location{StartColumn: 6, EndColumn: 11, Path_: testFile()},
					},
					IsClass:        true,
					Members_:       map[string]ast.Member{},
					Methods_:       map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
			}
			var d generics.Decl
			lookup := &type_.Lookup{
				Context_: &TypeContext{
					Project: &path.Project{},
					Scope:   []ast.Declaration{wrappingDecl(), c},
				},
				Tokens: []token.Token{
					{
						Type:  token.TOK_IDENTIFIER,
						Value: "T",
						Loc:   token.Location{StartColumn: 12, EndColumn: 13, Path_: testFile()},
					},
				},
			}
			d = generics.Decl{
				SymbolSlice: []ast.Type{
					&type_.Param{
						Context_: &TypeContext{
							Project: &path.Project{},
							Scope:   []ast.Declaration{wrappingDecl(), c},
							SymbolMap: SymbolMap{
								Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
							},
						},
						Token: lookup.Tokens[0],
						Decl:  &d,
					},
				},
				Declaration: c,
			}

			return testCase{
				name:     "valid empty generic class",
				input:    `class Class[T] {}`,
				expected: &d,
				errFunc:  assert.NoError,
			}
		}(),
		func() testCase {
			c := &decl.Class{
				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_CLASS,
						Loc:  token.Location{EndColumn: 5, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 17, EndColumn: 18, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Class",
						Loc:   token.Location{StartColumn: 6, EndColumn: 11, Path_: testFile()},
					},
					IsClass:        true,
					Members_:       map[string]ast.Member{},
					Methods_:       map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
				IsTailed_: true,
			}
			var d generics.Decl
			lookup := &type_.Lookup{
				Context_: &TypeContext{
					Project: &path.Project{},
					Scope:   []ast.Declaration{wrappingDecl(), c},
					SymbolMap: SymbolMap{
						Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
					},
				},
				Tokens: []token.Token{
					{
						Type:  token.TOK_IDENTIFIER,
						Value: "T",
						Loc:   token.Location{StartColumn: 12, EndColumn: 13, Path_: testFile()},
					},
				},
			}
			d = generics.Decl{
				SymbolSlice: []ast.Type{
					&type_.Param{
						Context_: &TypeContext{
							Project: &path.Project{},
							Scope:   []ast.Declaration{wrappingDecl(), c},
							SymbolMap: SymbolMap{
								Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
							},
						},
						Token: lookup.Tokens[0],
						Decl:  &d,
					},
				},
				Declaration: c,
			}
			return testCase{
				name:     "valid empty tailed generic class",
				input:    `class Class[T]~ {}`,
				expected: &d,
				errFunc:  assert.NoError,
			}
		}(),
		{
			name:  "valid empty abstract",
			input: `abstract Abstract {}`,
			expected: &decl.Abstract{
				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_ABSTRACT,
						Loc:  token.Location{EndColumn: 8, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 19, EndColumn: 20, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Abstract",
						Loc:   token.Location{StartColumn: 9, EndColumn: 17, Path_: testFile()},
					},
					Members_:       map[string]ast.Member{},
					Methods_:       map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty interface",
			input: `interface Interface {}`,
			expected: &decl.Interface{
				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_INTERFACE,
						Loc:  token.Location{EndColumn: 9, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 21, EndColumn: 22, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Interface",
						Loc:   token.Location{StartColumn: 10, EndColumn: 19, Path_: testFile()},
					},
					Members_:       map[string]ast.Member{},
					Methods_:       map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty struct",
			input: `struct Struct {}`,
			expected: &decl.Struct{
				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_STRUCT,
						Loc:  token.Location{EndColumn: 6, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 15, EndColumn: 16, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Struct",
						Loc:   token.Location{StartColumn: 7, EndColumn: 13, Path_: testFile()},
					},
					Members_:       map[string]ast.Member{},
					Methods_:       map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
			},
			errFunc: assert.NoError,
		},
		func() testCase {
			var enum decl.Enum
			enum = decl.Enum{
				BaseDecl: decl.BaseDecl{
					StartToken: token.Token{
						Type: token.TOK_ENUM,
						Loc:  token.Location{EndColumn: 4, Path_: testFile()},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc:  token.Location{StartColumn: 15, EndColumn: 16, Path_: testFile()},
					},
					Name_: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "Enum",
						Loc:   token.Location{StartColumn: 5, EndColumn: 9, Path_: testFile()},
					},
					Members_: map[string]ast.Member{},
					Methods_: map[string]ast.Method{},
					StaticMembers_: map[string]ast.Member{
						"ONE": &field.Enum{
							StartToken: token.Token{
								Type:  token.TOK_IDENTIFIER,
								Value: "ONE",
								Loc:   token.Location{StartColumn: 11, EndColumn: 14, Path_: testFile()},
							},
							EndToken: token.Token{
								Type: token.TOK_SEMICOLON,
								Loc:  token.Location{StartColumn: 14, EndColumn: 15, Path_: testFile()},
							},
							Type_: &enum,
							Expression: &expr.Identifier{
								Token: token.Token{
									Type:  token.TOK_IDENTIFIER,
									Value: "ONE",
									Loc:   token.Location{StartColumn: 11, EndColumn: 14, Path_: testFile()},
								},
							},
						},
					},
					StaticMethods_: map[string]ast.Method{},
					Declarations_:  map[string]ast.DeclField{},
				},
			}

			return testCase{
				name:     "valid enum",
				input:    `enum Enum {ONE;}`,
				expected: &enum,
				errFunc:  assert.NoError,
			}
		}(),
	}

	testParseHelper(t, testCases,
		func(p *Parser) (ast.Node, io.Error) {
			return p.ParseDecl()
		},
	)
}

func TestParseField(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name:  "valid empty member",
			input: `var int i;`,
			expected: &field.Member{
				Modifiers: field.Modifiers{},
				StartToken: token.Token{
					Type: token.TOK_VAR,
					Loc: token.Location{
						EndColumn: 3,
						Path_:     testFile(),
					},
				},
				EndToken: token.Token{
					Type: token.TOK_SEMICOLON,
					Loc: token.Location{
						StartColumn: 9,
						EndColumn:   10,
						Path_:       testFile(),
					},
				},
				Type_: &type_.Integer{
					Primitive: type_.Primitive[type_.Integer]{
						Token: token.Token{
							Type: token.TOK_TYPEINT,
							Loc: token.Location{
								StartColumn: 4,
								EndColumn:   7,
								Path_:       testFile(),
							},
						},
					},
				},
				Name_: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "i",
					Loc: token.Location{
						StartColumn: 8,
						EndColumn:   9,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid nonempty member",
			input: `var int i = 1;`,
			expected: &field.Member{
				Modifiers: field.Modifiers{},
				StartToken: token.Token{
					Type: token.TOK_VAR,
					Loc: token.Location{
						EndColumn: 3,
						Path_:     testFile(),
					},
				},
				EndToken: token.Token{
					Type: token.TOK_SEMICOLON,
					Loc: token.Location{
						StartColumn: 13,
						EndColumn:   14,
						Path_:       testFile(),
					},
				},
				Type_: &type_.Integer{
					Primitive: type_.Primitive[type_.Integer]{
						Token: token.Token{
							Type: token.TOK_TYPEINT,
							Loc: token.Location{
								StartColumn: 4,
								EndColumn:   7,
								Path_:       testFile(),
							},
						},
					},
				},
				Name_: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "i",
					Loc: token.Location{
						StartColumn: 8,
						EndColumn:   9,
						Path_:       testFile(),
					},
				},
				Expr: &expr.Integer{
					Token: token.Token{
						Type:    token.TOK_INTEGER,
						Value:   "1",
						Integer: 1,
						Loc: token.Location{
							StartColumn: 12,
							EndColumn:   13,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty method",
			input: `fun void() f;`,
			expected: &field.Method{
				Modifiers: field.Modifiers{},
				StartToken: token.Token{
					Type: token.TOK_FUN,
					Loc: token.Location{
						EndColumn: 3,
						Path_:     testFile(),
					},
				},
				EndToken: token.Token{
					Type: token.TOK_SEMICOLON,
					Loc: token.Location{
						StartColumn: 12,
						EndColumn:   13,
						Path_:       testFile(),
					},
				},
				Name_: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "f",
					Loc: token.Location{
						StartColumn: 11,
						EndColumn:   12,
						Path_:       testFile(),
					},
				},
				Type_: &type_.Function{
					Constant: true,
					ReturnType_: &type_.Void{
						Primitive: type_.Primitive[type_.Void]{
							Token: token.Token{
								Type: token.TOK_TYPEVOID,
								Loc: token.Location{
									StartColumn: 4,
									EndColumn:   8,
									Path_:       testFile(),
								},
							},
						},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTPAREN,
						Loc: token.Location{
							StartColumn: 9,
							EndColumn:   10,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid nonempty method",
			input: `fun void() f = (){}`,
			expected: &field.Method{
				Modifiers: field.Modifiers{},
				StartToken: token.Token{
					Type: token.TOK_FUN,
					Loc: token.Location{
						EndColumn: 3,
						Path_:     testFile(),
					},
				},
				Name_: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "f",
					Loc: token.Location{
						StartColumn: 11,
						EndColumn:   12,
						Path_:       testFile(),
					},
				},
				Type_: &type_.Function{
					Constant: true,
					ReturnType_: &type_.Void{
						Primitive: type_.Primitive[type_.Void]{
							Token: token.Token{
								Type: token.TOK_TYPEVOID,
								Loc: token.Location{
									StartColumn: 4,
									EndColumn:   8,
									Path_:       testFile(),
								},
							},
						},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTPAREN,
						Loc: token.Location{
							StartColumn: 9,
							EndColumn:   10,
							Path_:       testFile(),
						},
					},
				},
				Stmt: &stmt.Block{
					BoundedStmt: stmt.BoundedStmt{
						StartToken: token.Token{
							Type: token.TOK_LEFTBRACE,
							Loc: token.Location{
								StartColumn: 17,
								EndColumn:   18,
								Path_:       testFile(),
							},
						},
						EndToken: token.Token{
							Type: token.TOK_RIGHTBRACE,
							Loc: token.Location{
								StartColumn: 18,
								EndColumn:   19,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
	}

	testParseHelper(t, testCases,
		func(p *Parser) (ast.Node, io.Error) {
			return p.ParseField()
		},
	)
}

func TestParseStmt(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name:  "valid empty block",
			input: `{}`,
			expected: &stmt.Block{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_LEFTBRACE,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   1,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc: token.Location{
							StartColumn: 1,
							EndColumn:   2,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty break",
			input: `break;`,
			expected: &stmt.Break{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_BREAK,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   5,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 5,
							EndColumn:   6,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid labelled break",
			input: `break b;`,
			expected: &stmt.Break{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_BREAK,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   5,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 7,
							EndColumn:   8,
							Path_:       testFile(),
						},
					},
				},
				Label: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "b",
					Loc: token.Location{
						StartColumn: 6,
						EndColumn:   7,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty continue",
			input: `continue;`,
			expected: &stmt.Continue{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_CONTINUE,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   8,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 8,
							EndColumn:   9,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid labelled continue",
			input: `continue c;`,
			expected: &stmt.Continue{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_CONTINUE,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   8,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 10,
							EndColumn:   11,
							Path_:       testFile(),
						},
					},
				},
				Label: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "c",
					Loc: token.Location{
						StartColumn: 9,
						EndColumn:   10,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid delete",
			input: `delete d;`,
			expected: &stmt.Delete{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_DELETE,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   6,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 8,
							EndColumn:   9,
							Path_:       testFile(),
						},
					},
				},
				Expr: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "d",
						Loc: token.Location{
							StartColumn: 7,
							EndColumn:   8,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid expr statement",
			input: `e;`,
			expected: &stmt.Expr{
				Expr: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "e",
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   1,
							Path_:       testFile(),
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_SEMICOLON,
					Loc: token.Location{
						StartColumn: 1,
						EndColumn:   2,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid for0",
			input: `for {}`,
			expected: &stmt.For0{
				StartToken: token.Token{
					Type: token.TOK_FOR,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   3,
						Path_:       testFile(),
					},
				},
				Stmt: &stmt.Block{
					BoundedStmt: stmt.BoundedStmt{
						StartToken: token.Token{
							Type: token.TOK_LEFTBRACE,
							Loc: token.Location{
								StartColumn: 4,
								EndColumn:   5,
								Path_:       testFile(),
							},
						},
						EndToken: token.Token{
							Type: token.TOK_RIGHTBRACE,
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   6,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid for1",
			input: `for (b) {} else {}`,
			expected: &stmt.For1{
				StartToken: token.Token{
					Type: token.TOK_FOR,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   3,
						Path_:       testFile(),
					},
				},
				Condition: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "b",
						Loc: token.Location{
							StartColumn: 5,
							EndColumn:   6,
							Path_:       testFile(),
						},
					},
				},
				Then: &stmt.Block{
					BoundedStmt: stmt.BoundedStmt{
						StartToken: token.Token{
							Type: token.TOK_LEFTBRACE,
							Loc: token.Location{
								StartColumn: 8,
								EndColumn:   9,
								Path_:       testFile(),
							},
						},
						EndToken: token.Token{
							Type: token.TOK_RIGHTBRACE,
							Loc: token.Location{
								StartColumn: 9,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
				Else: &stmt.Block{
					BoundedStmt: stmt.BoundedStmt{
						StartToken: token.Token{
							Type: token.TOK_LEFTBRACE,
							Loc: token.Location{
								StartColumn: 16,
								EndColumn:   17,
								Path_:       testFile(),
							},
						},
						EndToken: token.Token{
							Type: token.TOK_RIGHTBRACE,
							Loc: token.Location{
								StartColumn: 17,
								EndColumn:   18,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid for3",
			input: `for (var int i = 0; b; i++) {} else {}`,
			expected: &stmt.Block{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_FOR,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   3,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACE,
						Loc: token.Location{
							StartColumn: 37,
							EndColumn:   38,
							Path_:       testFile(),
						},
					},
				},
				Stmts: []ast.Statement{
					&stmt.Var{
						BoundedStmt: stmt.BoundedStmt{
							StartToken: token.Token{
								Type: token.TOK_VAR,
								Loc: token.Location{
									StartColumn: 5,
									EndColumn:   8,
									Path_:       testFile(),
								},
							},
							EndToken: token.Token{
								Type: token.TOK_SEMICOLON,
								Loc: token.Location{
									StartColumn: 18,
									EndColumn:   19,
									Path_:       testFile(),
								},
							},
						},
						Type_: &type_.Integer{
							Primitive: type_.Primitive[type_.Integer]{
								Token: token.Token{
									Type: token.TOK_TYPEINT,
									Loc: token.Location{
										StartColumn: 9,
										EndColumn:   12,
										Path_:       testFile(),
									},
								},
							},
						},
						Name_: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 13,
								EndColumn:   14,
								Path_:       testFile(),
							},
						},
						Expr: &expr.Integer{
							Token: token.Token{
								Type:  token.TOK_INTEGER,
								Value: "0",
								Loc: token.Location{
									StartColumn: 17,
									EndColumn:   18,
									Path_:       testFile(),
								},
							},
						},
					},
					&stmt.For1{
						StartToken: token.Token{
							Type: token.TOK_FOR,
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   3,
								Path_:       testFile(),
							},
						},
						Condition: &expr.Identifier{
							Token: token.Token{
								Type:  token.TOK_IDENTIFIER,
								Value: "b",
								Loc: token.Location{
									StartColumn: 20,
									EndColumn:   21,
									Path_:       testFile(),
								},
							},
						},
						Then: &stmt.Block{
							BoundedStmt: stmt.BoundedStmt{
								StartToken: token.Token{
									Type: token.TOK_LEFTBRACE,
									Loc: token.Location{
										StartColumn: 28,
										EndColumn:   29,
										Path_:       testFile(),
									},
								},
								EndToken: token.Token{
									Type: token.TOK_RIGHTBRACE,
									Loc: token.Location{
										StartColumn: 37,
										EndColumn:   38,
										Path_:       testFile(),
									},
								},
							},
							Stmts: []ast.Statement{
								&stmt.Block{
									BoundedStmt: stmt.BoundedStmt{
										StartToken: token.Token{
											Type: token.TOK_LEFTBRACE,
											Loc: token.Location{
												StartColumn: 28,
												EndColumn:   29,
												Path_:       testFile(),
											},
										},
										EndToken: token.Token{
											Type: token.TOK_RIGHTBRACE,
											Loc: token.Location{
												StartColumn: 29,
												EndColumn:   30,
												Path_:       testFile(),
											},
										},
									},
								},
								&stmt.Expr{
									Expr: &expr.IncrementSuffix{
										SuffixedUnary: expr.SuffixedUnary{
											Token: token.Token{
												Type: token.TOK_INC,
												Loc: token.Location{
													StartColumn: 24,
													EndColumn:   26,
													Path_:       testFile(),
												},
											},
											Expr: &expr.Identifier{
												Token: token.Token{
													Type:  token.TOK_IDENTIFIER,
													Value: "i",
													Loc: token.Location{
														StartColumn: 23,
														EndColumn:   24,
														Path_:       testFile(),
													},
												},
											},
										},
									},
									EndToken: token.Token{
										Type: token.TOK_INC,
										Loc: token.Location{
											StartColumn: 24,
											EndColumn:   26,
											Path_:       testFile(),
										},
									},
								},
							},
						},
						Else: &stmt.Block{
							BoundedStmt: stmt.BoundedStmt{
								StartToken: token.Token{
									Type: token.TOK_LEFTBRACE,
									Loc: token.Location{
										StartColumn: 36,
										EndColumn:   37,
										Path_:       testFile(),
									},
								},
								EndToken: token.Token{
									Type: token.TOK_RIGHTBRACE,
									Loc: token.Location{
										StartColumn: 37,
										EndColumn:   38,
										Path_:       testFile(),
									},
								},
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid if",
			input: `if (b) {} else {}`,
			expected: &stmt.If{
				StartToken: token.Token{
					Type: token.TOK_IF,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   2,
						Path_:       testFile(),
					},
				},
				Condition: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "b",
						Loc: token.Location{
							StartColumn: 4,
							EndColumn:   5,
							Path_:       testFile(),
						},
					},
				},
				Then: &stmt.Block{
					BoundedStmt: stmt.BoundedStmt{
						StartToken: token.Token{
							Type: token.TOK_LEFTBRACE,
							Loc: token.Location{
								StartColumn: 7,
								EndColumn:   8,
								Path_:       testFile(),
							},
						},
						EndToken: token.Token{
							Type: token.TOK_RIGHTBRACE,
							Loc: token.Location{
								StartColumn: 8,
								EndColumn:   9,
								Path_:       testFile(),
							},
						},
					},
				},
				Else: &stmt.Block{
					BoundedStmt: stmt.BoundedStmt{
						StartToken: token.Token{
							Type: token.TOK_LEFTBRACE,
							Loc: token.Location{
								StartColumn: 15,
								EndColumn:   16,
								Path_:       testFile(),
							},
						},
						EndToken: token.Token{
							Type: token.TOK_RIGHTBRACE,
							Loc: token.Location{
								StartColumn: 16,
								EndColumn:   17,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid label",
			input: `label l : for {}`,
			expected: &stmt.Label{
				StartToken: token.Token{
					Type: token.TOK_LABEL,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   5,
						Path_:       testFile(),
					},
				},
				Label: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "l",
					Loc: token.Location{
						StartColumn: 6,
						EndColumn:   7,
						Path_:       testFile(),
					},
				},
				Stmt: &stmt.For0{
					StartToken: token.Token{
						Type: token.TOK_FOR,
						Loc: token.Location{
							StartColumn: 10,
							EndColumn:   13,
							Path_:       testFile(),
						},
					},
					Stmt: &stmt.Block{
						BoundedStmt: stmt.BoundedStmt{
							StartToken: token.Token{
								Type: token.TOK_LEFTBRACE,
								Loc: token.Location{
									StartColumn: 14,
									EndColumn:   15,
									Path_:       testFile(),
								},
							},
							EndToken: token.Token{
								Type: token.TOK_RIGHTBRACE,
								Loc: token.Location{
									StartColumn: 15,
									EndColumn:   16,
									Path_:       testFile(),
								},
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid panic",
			input: `panic("panic");`,
			expected: &stmt.Panic{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_PANIC,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   5,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 14,
							EndColumn:   15,
							Path_:       testFile(),
						},
					},
				},
				Expr: &expr.String{
					Token: token.Token{
						Type:  token.TOK_STRING,
						Value: "panic",
						Loc: token.Location{
							StartColumn: 6,
							EndColumn:   13,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid print",
			input: `print("print");`,
			expected: &stmt.Print{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_PRINT,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   5,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 14,
							EndColumn:   15,
							Path_:       testFile(),
						},
					},
				},
				Expr: &expr.String{
					Token: token.Token{
						Type:  token.TOK_STRING,
						Value: "print",
						Loc: token.Location{
							StartColumn: 6,
							EndColumn:   13,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty return",
			input: `return;`,
			expected: &stmt.Return{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_RETURN,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   6,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 6,
							EndColumn:   7,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid empty return",
			input: `return r;`,
			expected: &stmt.Return{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_RETURN,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   6,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 8,
							EndColumn:   9,
							Path_:       testFile(),
						},
					},
				},
				Expr: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "r",
						Loc: token.Location{
							StartColumn: 7,
							EndColumn:   8,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "valid var",
			input: `var int i = 0;`,
			expected: &stmt.Var{
				BoundedStmt: stmt.BoundedStmt{
					StartToken: token.Token{
						Type: token.TOK_VAR,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   3,
							Path_:       testFile(),
						},
					},
					EndToken: token.Token{
						Type: token.TOK_SEMICOLON,
						Loc: token.Location{
							StartColumn: 13,
							EndColumn:   14,
							Path_:       testFile(),
						},
					},
				},
				Type_: &type_.Integer{
					Primitive: type_.Primitive[type_.Integer]{
						Token: token.Token{
							Type: token.TOK_TYPEINT,
							Loc: token.Location{
								StartColumn: 4,
								EndColumn:   7,
								Path_:       testFile(),
							},
						},
					},
				},
				Name_: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "i",
					Loc: token.Location{
						StartColumn: 8,
						EndColumn:   9,
						Path_:       testFile(),
					},
				},
				Expr: &expr.Integer{
					Token: token.Token{
						Type:  token.TOK_INTEGER,
						Value: "0",
						Loc: token.Location{
							StartColumn: 12,
							EndColumn:   13,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
	}

	testParseHelper(t, testCases,
		func(p *Parser) (ast.Node, io.Error) {
			return p.ParseStmt()
		},
	)
}

func TestParseExpr(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name:  "access",
			input: `left[right]`,
			expected: &expr.Access{
				Left: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "left",
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   4,
							Path_:       testFile(),
						},
					},
				},
				Right: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "right",
						Loc: token.Location{
							StartColumn: 5,
							EndColumn:   10,
							Path_:       testFile(),
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_RIGHTBRACKET,
					Loc: token.Location{
						StartColumn: 10,
						EndColumn:   11,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "access",
			input: `left=right`,
			expected: &expr.Assign{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "async",
			input: `async f()`,
			expected: &expr.Async{
				StartToken: token.Token{
					Type: token.TOK_ASYNC,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   5,
						Path_:       testFile(),
					},
				},
				CallExpr: &expr.Call{
					SuffixedToken: expr.SuffixedToken{
						Expr: &expr.Identifier{
							Token: token.Token{
								Type:  token.TOK_IDENTIFIER,
								Value: "f",
								Loc: token.Location{
									StartColumn: 6,
									EndColumn:   7,
									Path_:       testFile(),
								},
							},
						},
						Token: token.Token{
							Type: token.TOK_RIGHTPAREN,
							Loc: token.Location{
								StartColumn: 8,
								EndColumn:   9,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "bitwise and",
			input: `left&right`,
			expected: &expr.BitwiseAnd{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "bitwise or",
			input: `left|right`,
			expected: &expr.BitwiseOr{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "bitwise xor",
			input: `left^right`,
			expected: &expr.BitwiseXor{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "call",
			input: `f(a)`,
			expected: &expr.Call{
				SuffixedToken: expr.SuffixedToken{
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "f",
							Loc: token.Location{
								EndColumn: 1,
								Path_:     testFile(),
							},
						},
					},
					Token: token.Token{
						Type: token.TOK_RIGHTPAREN,
						Loc: token.Location{
							StartColumn: 3,
							EndColumn:   4,
							Path_:       testFile(),
						},
					},
				},
				Exprs: []ast.Expression{
					&expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "a",
							Loc: token.Location{
								StartColumn: 2,
								EndColumn:   3,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "less than",
			input: `left<right`,
			expected: &expr.LessThan{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "less than or equal",
			input: `left<=right`,
			expected: &expr.LessThanOrEqual{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "greater than",
			input: `left>right`,
			expected: &expr.GreaterThan{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "greater than or equal",
			input: `left>=right`,
			expected: &expr.GreaterThanOrEqual{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "spaceship",
			input: `left<=>right`,
			expected: &expr.Spaceship{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 7,
								EndColumn:   12,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "equal",
			input: `left==right`,
			expected: &expr.Equal{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "unequal",
			input: `left!=right`,
			expected: &expr.BangEqual{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "multiply",
			input: `left*right`,
			expected: &expr.Multiply{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "divide",
			input: `left/right`,
			expected: &expr.Divide{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "modulo",
			input: `left%right`,
			expected: &expr.Modulo{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "field",
			input: `left.right`,
			expected: &expr.Field{
				SuffixedToken: expr.SuffixedToken{
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "right",
						Loc: token.Location{
							StartColumn: 5,
							EndColumn:   10,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "type field",
			input: `type(left).right`,
			expected: &expr.TypeField{
				StartToken: token.Token{
					Type: token.TOK_TYPE,
					Loc: token.Location{
						EndColumn: 4,
						Path_:     testFile(),
					},
				},
				Type: &type_.Lookup{
					Context_: &TypeContext{
						Project: &path.Project{},
						Scope:   []ast.Declaration{wrappingDecl()},
						SymbolMap: SymbolMap{
							Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
						},
					},
					Tokens: []token.Token{
						{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   9,
								Path_:       testFile(),
							},
						},
					},
				},
				Token: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "right",
					Loc: token.Location{
						StartColumn: 11,
						EndColumn:   16,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "hash",
			input: `#left`,
			expected: &expr.Hash{
				PrefixedToken: expr.PrefixedToken{
					Token: token.Token{
						Type: token.TOK_HASHTAG,
						Loc: token.Location{
							EndColumn: 1,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 1,
								EndColumn:   5,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "hash",
			input: `i`,
			expected: &expr.Identifier{
				Token: token.Token{
					Type:  token.TOK_IDENTIFIER,
					Value: "i",
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   1,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "negate",
			input: `-i`,
			expected: &expr.Negation{
				PrefixedUnary: expr.PrefixedUnary{
					Token: token.Token{
						Type: token.TOK_MINUS,
						Loc: token.Location{
							EndColumn: 1,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 1,
								EndColumn:   2,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "bang",
			input: `!b`,
			expected: &expr.Bang{
				PrefixedUnary: expr.PrefixedUnary{
					Token: token.Token{
						Type: token.TOK_BANG,
						Loc: token.Location{
							EndColumn: 1,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "b",
							Loc: token.Location{
								StartColumn: 1,
								EndColumn:   2,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "lambda",
			input: `lambda void() : () {}`,
			expected: &expr.Lambda{
				StartToken: token.Token{
					Type: token.TOK_LAMBDA,
					Loc: token.Location{
						EndColumn: 6,
						Path_:     testFile(),
					},
				},
				Type: &type_.Function{
					ReturnType_: &type_.Void{
						Primitive: type_.Primitive[type_.Void]{
							Token: token.Token{
								Type: token.TOK_TYPEVOID,
								Loc: token.Location{
									StartColumn: 7,
									EndColumn:   11,
									Path_:       testFile(),
								},
							},
						},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTPAREN,
						Loc: token.Location{
							StartColumn: 12,
							EndColumn:   13,
							Path_:       testFile(),
						},
					},
				},
				Stmt: &stmt.Block{
					BoundedStmt: stmt.BoundedStmt{
						StartToken: token.Token{
							Type: token.TOK_LEFTBRACE,
							Loc: token.Location{
								StartColumn: 19,
								EndColumn:   20,
								Path_:       testFile(),
							},
						},
						EndToken: token.Token{
							Type: token.TOK_RIGHTBRACE,
							Loc: token.Location{
								StartColumn: 20,
								EndColumn:   21,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "and",
			input: `left&&right`,
			expected: &expr.And{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "or",
			input: `left||right`,
			expected: &expr.Or{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "inc pre",
			input: `++i`,
			expected: &expr.IncrementPrefix{
				PrefixedUnary: expr.PrefixedUnary{
					Token: token.Token{
						Type: token.TOK_INC,
						Loc: token.Location{
							EndColumn: 2,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 2,
								EndColumn:   3,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "dec pre",
			input: `--i`,
			expected: &expr.DecrementPrefix{
				PrefixedUnary: expr.PrefixedUnary{
					Token: token.Token{
						Type: token.TOK_DEC,
						Loc: token.Location{
							EndColumn: 2,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 2,
								EndColumn:   3,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "inc post",
			input: `i++`,
			expected: &expr.IncrementSuffix{
				SuffixedUnary: expr.SuffixedUnary{
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								EndColumn: 1,
								Path_:     testFile(),
							},
						},
					},
					Token: token.Token{
						Type: token.TOK_INC,
						Loc: token.Location{
							StartColumn: 1,
							EndColumn:   3,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "dec post",
			input: `i--`,
			expected: &expr.DecrementSuffix{
				SuffixedUnary: expr.SuffixedUnary{
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								EndColumn: 1,
								Path_:     testFile(),
							},
						},
					},
					Token: token.Token{
						Type: token.TOK_DEC,
						Loc: token.Location{
							StartColumn: 1,
							EndColumn:   3,
							Path_:       testFile(),
						},
					}},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "new",
			input: `new(Node, 2, 3)`,
			expected: &expr.New{
				StartToken: token.Token{
					Type: token.TOK_NEW,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   3,
						Path_:       testFile(),
					},
				},
				EndToken: token.Token{
					Type: token.TOK_RIGHTPAREN,
					Loc: token.Location{
						StartColumn: 14,
						EndColumn:   15,
						Path_:       testFile(),
					},
				},
				Type: &type_.Lookup{
					Context_: &TypeContext{
						Project: &path.Project{},
						Scope:   []ast.Declaration{wrappingDecl()},
						SymbolMap: SymbolMap{
							Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
						},
					},
					Tokens: []token.Token{
						{
							Type:  token.TOK_IDENTIFIER,
							Value: "Node",
							Loc: token.Location{
								StartColumn: 4,
								EndColumn:   8,
								Path_:       testFile(),
							},
						},
					},
				},
				TailLengthExpr: &expr.Integer{
					Token: token.Token{
						Type:    token.TOK_INTEGER,
						Integer: 2,
						Value:   "2",
						Loc: token.Location{
							StartColumn: 10,
							EndColumn:   11,
							Path_:       testFile(),
						},
					},
				},
				ArrayLengthExpr: &expr.Integer{
					Token: token.Token{
						Type:    token.TOK_INTEGER,
						Integer: 3,
						Value:   "3",
						Loc: token.Location{
							StartColumn: 13,
							EndColumn:   14,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "nullary",
			input: `left??right`,
			expected: &expr.Nullary{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "this",
			input: `this`,
			expected: &expr.This{
				Token: token.Token{
					Type: token.TOK_THIS,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   4,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "super",
			input: `super`,
			expected: &expr.Super{
				Token: token.Token{
					Type: token.TOK_SUPER,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   5,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "dereference",
			input: `*i`,
			expected: &expr.Dereference{
				PrefixedToken: expr.PrefixedToken{
					Token: token.Token{
						Type: token.TOK_STAR,
						Loc: token.Location{
							EndColumn: 1,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 1,
								EndColumn:   2,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "reference",
			input: `&i`,
			expected: &expr.Reference{
				PrefixedToken: expr.PrefixedToken{
					Token: token.Token{
						Type: token.TOK_BITAND,
						Loc: token.Location{
							EndColumn: 1,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 1,
								EndColumn:   2,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "integer",
			input: `1`,
			expected: &expr.Integer{
				Token: token.Token{
					Type:    token.TOK_INTEGER,
					Integer: 1,
					Value:   "1",
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   1,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "float",
			input: `1.0`,
			expected: &expr.Float{
				Token: token.Token{
					Type:  token.TOK_FLOAT,
					Float: 1,
					Value: "1.0",
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   3,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "character",
			input: `'c'`,
			expected: &expr.Character{
				Token: token.Token{
					Type:  token.TOK_CHARACTER,
					Rune:  'c',
					Value: "c",
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   3,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "string",
			input: `"s"`,
			expected: &expr.String{
				Token: token.Token{
					Type:  token.TOK_STRING,
					Value: `s`,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   3,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "true",
			input: `true`,
			expected: &expr.True{
				Token: token.Token{
					Type: token.TOK_TRUE,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   4,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "false",
			input: `false`,
			expected: &expr.False{
				Token: token.Token{
					Type: token.TOK_FALSE,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   5,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "null",
			input: `null`,
			expected: &expr.Null{
				Token: token.Token{
					Type: token.TOK_NULL,
					Loc: token.Location{
						StartColumn: 0,
						EndColumn:   4,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "left shift",
			input: `left<<right`,
			expected: &expr.LeftShift{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "right shift",
			input: `left>>right`,
			expected: &expr.RightShift{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 6,
								EndColumn:   11,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "unsigned right shift",
			input: `left>>>right`,
			expected: &expr.UnsignedRightShift{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 7,
								EndColumn:   12,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "add",
			input: `left+right`,
			expected: &expr.Add{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "subtract",
			input: `left-right`,
			expected: &expr.Subtract{
				Binary: expr.Binary{
					Left: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "left",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
					Right: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "right",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   10,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "ternary",
			input: `cond?left:right`,
			expected: &expr.Ternary{
				Condition: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "cond",
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   4,
							Path_:       testFile(),
						},
					},
				},
				Then: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "left",
						Loc: token.Location{
							StartColumn: 5,
							EndColumn:   9,
							Path_:       testFile(),
						},
					},
				},
				Else: &expr.Identifier{
					Token: token.Token{
						Type:  token.TOK_IDENTIFIER,
						Value: "right",
						Loc: token.Location{
							StartColumn: 10,
							EndColumn:   15,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "type of",
			input: `@i`,
			expected: &expr.TypeOf{
				PrefixedToken: expr.PrefixedToken{
					Token: token.Token{
						Type: token.TOK_AT,
						Loc: token.Location{
							EndColumn: 1,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 1,
								EndColumn:   2,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "wait",
			input: `wait i`,
			expected: &expr.Wait{
				PrefixedToken: expr.PrefixedToken{
					Token: token.Token{
						Type: token.TOK_WAIT,
						Loc: token.Location{
							EndColumn: 4,
							Path_:     testFile(),
						},
					},
					Expr: &expr.Identifier{
						Token: token.Token{
							Type:  token.TOK_IDENTIFIER,
							Value: "i",
							Loc: token.Location{
								StartColumn: 5,
								EndColumn:   6,
								Path_:       testFile(),
							},
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
	}

	testParseHelper(t, testCases,
		func(p *Parser) (ast.Node, io.Error) {
			return p.ParseExpr()
		},
	)
}

func TestParseType(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name:  "array",
			input: `int[8]`,
			expected: &type_.Array{
				Type: &type_.Integer{
					Primitive: type_.Primitive[type_.Integer]{
						Token: token.Token{
							Type: token.TOK_TYPEINT,
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   3,
								Path_:       testFile(),
							},
						},
					},
				},
				Size: 8,
				EndToken: token.Token{
					Type: token.TOK_RIGHTBRACKET,
					Loc: token.Location{
						StartColumn: 5,
						EndColumn:   6,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "multidimenstional array",
			input: `int[8][8]`,
			expected: &type_.Array{
				Type: &type_.Array{
					Type: &type_.Integer{
						Primitive: type_.Primitive[type_.Integer]{
							Token: token.Token{
								Type: token.TOK_TYPEINT,
								Loc: token.Location{
									StartColumn: 0,
									EndColumn:   3,
									Path_:       testFile(),
								},
							},
						},
					},
					Size: 8,
					EndToken: token.Token{
						Type: token.TOK_RIGHTBRACKET,
						Loc: token.Location{
							StartColumn: 5,
							EndColumn:   6,
							Path_:       testFile(),
						},
					},
				},
				Size: 8,
				EndToken: token.Token{
					Type: token.TOK_RIGHTBRACKET,
					Loc: token.Location{
						StartColumn: 8,
						EndColumn:   9,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "function",
			input: `void()`,
			expected: &type_.Function{
				ReturnType_: &type_.Void{
					Primitive: type_.Primitive[type_.Void]{
						Token: token.Token{
							Type: token.TOK_TYPEVOID,
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_RIGHTPAREN,
					Loc: token.Location{
						StartColumn: 5,
						EndColumn:   6,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "function with parameters",
			input: `void(int, int)`,
			expected: &type_.Function{
				ReturnType_: &type_.Void{
					Primitive: type_.Primitive[type_.Void]{
						Token: token.Token{
							Type: token.TOK_TYPEVOID,
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
				},
				ParameterTypes_: []ast.Type{
					&type_.Integer{
						Primitive: type_.Primitive[type_.Integer]{
							Token: token.Token{
								Type: token.TOK_TYPEINT,
								Loc: token.Location{
									StartColumn: 5,
									EndColumn:   8,
									Path_:       testFile(),
								},
							},
						},
					},
					&type_.Integer{
						Primitive: type_.Primitive[type_.Integer]{
							Token: token.Token{
								Type: token.TOK_TYPEINT,
								Loc: token.Location{
									StartColumn: 10,
									EndColumn:   13,
									Path_:       testFile(),
								},
							},
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_RIGHTPAREN,
					Loc: token.Location{
						StartColumn: 13,
						EndColumn:   14,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "function returned by function",
			input: `void()()`,
			expected: &type_.Function{
				ReturnType_: &type_.Function{
					ReturnType_: &type_.Void{
						Primitive: type_.Primitive[type_.Void]{
							Token: token.Token{
								Type: token.TOK_TYPEVOID,
								Loc: token.Location{
									StartColumn: 0,
									EndColumn:   4,
									Path_:       testFile(),
								},
							},
						},
					},
					EndToken: token.Token{
						Type: token.TOK_RIGHTPAREN,
						Loc: token.Location{
							StartColumn: 5,
							EndColumn:   6,
							Path_:       testFile(),
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_RIGHTPAREN,
					Loc: token.Location{
						StartColumn: 7,
						EndColumn:   8,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "generic",
			input: `List[int]`,
			expected: &type_.Generic{
				Context_: &TypeContext{
					Project: &path.Project{},
					Scope:   []ast.Declaration{wrappingDecl()},
					SymbolMap: SymbolMap{
						Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
					},
				},
				Type: &type_.Lookup{
					Context_: &TypeContext{
						Project: &path.Project{},
						Scope:   []ast.Declaration{wrappingDecl()},
						SymbolMap: SymbolMap{
							Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
						},
					},
					Tokens: []token.Token{
						{
							Type:  token.TOK_IDENTIFIER,
							Value: "List",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
				},
				GenericTypes: []ast.Type{
					&type_.Integer{
						Primitive: type_.Primitive[type_.Integer]{
							Token: token.Token{
								Type: token.TOK_TYPEINT,
								Loc: token.Location{
									StartColumn: 5,
									EndColumn:   8,
									Path_:       testFile(),
								},
							},
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_RIGHTBRACKET,
					Loc: token.Location{
						StartColumn: 8,
						EndColumn:   9,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "generic, multiple types",
			input: `Map[str, int]`,
			expected: &type_.Generic{
				Context_: &TypeContext{
					Project: &path.Project{},
					Scope:   []ast.Declaration{wrappingDecl()},
					SymbolMap: SymbolMap{
						Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
					},
				},
				Type: &type_.Lookup{
					Context_: &TypeContext{
						Project: &path.Project{},
						Scope:   []ast.Declaration{wrappingDecl()},
						SymbolMap: SymbolMap{
							Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
						},
					},
					Tokens: []token.Token{
						{
							Type:  token.TOK_IDENTIFIER,
							Value: "Map",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   3,
								Path_:       testFile(),
							},
						},
					},
				},
				GenericTypes: []ast.Type{
					&type_.String{
						Primitive: type_.Primitive[type_.String]{
							Token: token.Token{
								Type: token.TOK_TYPESTR,
								Loc: token.Location{
									StartColumn: 4,
									EndColumn:   7,
									Path_:       testFile(),
								},
							},
						},
					},
					&type_.Integer{
						Primitive: type_.Primitive[type_.Integer]{
							Token: token.Token{
								Type: token.TOK_TYPEINT,
								Loc: token.Location{
									StartColumn: 9,
									EndColumn:   12,
									Path_:       testFile(),
								},
							},
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_RIGHTBRACKET,
					Loc: token.Location{
						StartColumn: 12,
						EndColumn:   13,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "pointer",
			input: `int*`,
			expected: &type_.Pointer{
				Type: &type_.Integer{
					Primitive: type_.Primitive[type_.Integer]{
						Token: token.Token{
							Type: token.TOK_TYPEINT,
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   3,
								Path_:       testFile(),
							},
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_STAR,
					Loc: token.Location{
						StartColumn: 3,
						EndColumn:   4,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "tailed",
			input: `Node~*`,
			expected: &type_.Pointer{
				Type: &type_.Tailed{
					Type: &type_.Lookup{
						Context_: &TypeContext{
							Project: &path.Project{},
							Scope:   []ast.Declaration{wrappingDecl()},
							SymbolMap: SymbolMap{
								Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
							},
						},
						Tokens: []token.Token{
							{
								Type:  token.TOK_IDENTIFIER,
								Value: "Node",
								Loc: token.Location{
									StartColumn: 0,
									EndColumn:   4,
									Path_:       testFile(),
								},
							},
						},
					},
					Size: -1,
					EndToken: token.Token{
						Type: token.TOK_TILDE,
						Loc: token.Location{
							StartColumn: 4,
							EndColumn:   5,
							Path_:       testFile(),
						},
					},
				},
				EndToken: token.Token{
					Type: token.TOK_STAR,
					Loc: token.Location{
						StartColumn: 5,
						EndColumn:   6,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "tailed with size",
			input: `Node~8`,
			expected: &type_.Tailed{
				Type: &type_.Lookup{
					Context_: &TypeContext{
						Project: &path.Project{},
						Scope:   []ast.Declaration{wrappingDecl()},
						SymbolMap: SymbolMap{
							Types: &typeid.Types{ListsMap: map[typeid.ListKey]uint32{}},
						},
					},
					Tokens: []token.Token{
						{
							Type:  token.TOK_IDENTIFIER,
							Value: "Node",
							Loc: token.Location{
								StartColumn: 0,
								EndColumn:   4,
								Path_:       testFile(),
							},
						},
					},
				},
				Size: 8,
				EndToken: token.Token{
					Type:    token.TOK_INTEGER,
					Integer: 8,
					Value:   "8",
					Loc: token.Location{
						StartColumn: 5,
						EndColumn:   6,
						Path_:       testFile(),
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:  "constant",
			input: `int$`,
			expected: &type_.Integer{
				Primitive: type_.Primitive[type_.Integer]{
					Constant: true,
					Token: token.Token{
						Type: token.TOK_TYPEINT,
						Loc: token.Location{
							StartColumn: 0,
							EndColumn:   3,
							Path_:       testFile(),
						},
					},
				},
			},
			errFunc: assert.NoError,
		},
		{
			name:    "double constant",
			input:   `int$$`,
			errFunc: assert.Error,
		},
	}

	testParseHelper(t, testCases,
		func(p *Parser) (ast.Node, io.Error) {
			return p.ParseType()
		},
	)
}
